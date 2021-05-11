package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/mocks"
)

type testDone struct {
	failed bool
	msg    string
}

var l zerolog.Logger

func TestMain(m *testing.M) {
	l = zerolog.New(nil)

	os.Exit(m.Run())
}

func TestNewMailer(t *testing.T) {
	done := make(chan testDone)

	m := NewMailer(&l, func(err error) {
		done <- testDone{
			failed: true,
			msg:    "errHandler: " + err.Error(),
		}
	})

	if m.queue == nil || m.errorsQueue == nil || m.done == nil {
		t.Errorf("want mailer with not nil channels got %v", m)
	}
	select {
	case result := <-done:
		t.Errorf(result.msg)
	default:
	}
}

func TestSend(t *testing.T) {
	m := Mailer{
		queue: make(chan *emailRequest),
	}
	subject := []byte("TestSend")
	data := []byte("message")
	m.Send([]string{mocks.ExampleUser.EmailAddress}, subject, data)

	select {
	case <-time.After(time.Millisecond * 50):
		t.Error("custom timeout exceeded: 50ms")
	case er := <-m.queue:
		if string(er.subject) != string(subject) {
			t.Errorf("want subject %q, got %q", subject, er.subject)
		}
		if string(er.data) != string(data) {
			t.Errorf("want data %q, got %q", data, er.data)
		}

		if len(er.recipients) != 1 || er.recipients[0] != mocks.ExampleUser.EmailAddress {
			t.Errorf("want one recipient: %s, got %v",
				mocks.ExampleUser.EmailAddress, er.recipients)
		}
	}

}

func TestClose(t *testing.T) {
	done := make(chan testDone)

	m := NewMailer(&l, func(err error) {
		done <- testDone{
			failed: true,
			msg:    "errHandler: " + err.Error(),
		}
	})

	go func() {
		m.Close()
		done <- testDone{}
	}()

	select {
	case <-time.After(time.Millisecond * 50):

		done <- testDone{
			failed: true,
			msg:    "custom timeout exceeded: 50ms",
		}
	case result := <-done:
		if result.failed {
			t.Error(result.msg)
		}
	}

}

func TestInternalSend(t *testing.T) {

	done := make(chan testDone)

	emailReq := &emailRequest{
		subject: []byte("--- A ---"),
		recipients: []string{
			mocks.ExampleUser.EmailAddress,
		},
		data: []byte(">> Data A <<"),
	}

	m := NewMailer(&l, func(err error) {
		done <- testDone{
			failed: true,
			msg:    "errHandler: " + err.Error(),
		}
	})
	m.sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		var errMsg []string
		if from != os.Getenv("APP_EMAIL_ADDRESS") {
			errMsg = append(errMsg, fmt.Sprintf("want from %s, got %s", os.Getenv("APP_EMAIL_ADDRESS"), from))
		}

		for i, e := range to {
			if emailReq.recipients[i] != e {
				errMsg = append(errMsg, fmt.Sprintf("want recipients %v, got %v", emailReq.recipients, to))
				break
			}
		}

		if !bytes.Contains(msg, emailReq.subject) {
			t.Errorf("want subject to be in msg payload")
		}

		if !bytes.Contains(msg, emailReq.data) {
			t.Errorf("want data to be in msg payload")
		}

		if len(errMsg) != 0 {
			msg := strings.Join(errMsg, "\n")
			done <- testDone{
				failed: true,
				msg:    msg,
			}

			return fmt.Errorf(msg)
		}

		done <- testDone{}
		return nil
	}

	go m.send(emailReq)

	select {
	case <-time.After(time.Millisecond * 50):
		t.Error("custom timeout exceeded: 50ms")
	case result := <-done:
		if result.failed {
			t.Error(result)
		}
	}

}

func TestErrorHandling(t *testing.T) {
	want := "test error message from"

	done := make(chan testDone)

	emailReq := &emailRequest{
		subject: []byte("--- A ---"),
		recipients: []string{
			mocks.ExampleUser.EmailAddress,
		},
		data: []byte(">> Data A <<"),
	}

	m := NewMailer(&l, func(err error) {
		if !strings.Contains(err.Error(), want) {
			done <- testDone{
				failed: true,
				msg:    "errHandler: " + err.Error(),
			}
		}
		done <- testDone{}
	})
	m.sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return fmt.Errorf(want)
	}

	go m.Send(emailReq.recipients, emailReq.subject, emailReq.data)

	select {
	case <-time.After(time.Millisecond * 50):
		t.Error("custom timeout exceeded: 50ms")
	case result := <-done:
		if result.failed {
			t.Errorf("want %q, got %q", want, result.msg)
		}
	}

}
