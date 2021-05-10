package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

var l = log.Default()

type Mailer struct {
	queue       chan *emailRequest
	done        chan struct{}
	errorsQueue chan error
	callback    ErrorHandler
}

type emailRequest struct {
	recipients []string
	data       []byte
}

// ErrorHandler is called when an error occurs, it guarantees that err is not nil
type ErrorHandler func(err error)

// NewMailer creates new Mailer, it setup listeners and returns created instace
func NewMailer(errHandler ErrorHandler) *Mailer {
	m := &Mailer{
		callback:    errHandler,
		queue:       make(chan *emailRequest, 5),
		errorsQueue: make(chan error, 5),
		done:        make(chan struct{}),
	}

	go m.listenErrors()
	go m.listen()
	return m
}

// Send emits event to send an email
func (m *Mailer) Send(recipients []string, data []byte) {

	select {
	case <-m.done:
		m.emitError(fmt.Errorf("1. mailer already closed"))
		return
	default:
		if len(recipients) == 0 {
			m.emitError(fmt.Errorf("missing mail recipients"))
			return
		}
		go func() {
			select {
			case <-m.done:
				m.emitError(fmt.Errorf("2. mailer already closed"))

			case m.queue <- &emailRequest{
				recipients: recipients,
				data:       data,
			}:
			}
		}()
	}
}

// Close stops the mailer service, it will cancell all queued tasks
// and waits for tasks in progress to complete before returning
func (m *Mailer) Close() {
	l.Println("[Close] about to close channels")
	close(m.done)

	l.Println("[Close] channels closed")
	wg := sync.WaitGroup{}
	wg.Add(2)
	// drain queues
	go func() {
		for er := range m.queue {
			l.Printf("[Close - drain] - message to: %s\n", er.recipients)
		}
		wg.Done()
	}()
	go func() {
		for err := range m.errorsQueue {
			l.Printf("[Close - drain] - error to: %s\n", err)
		}
		wg.Done()
	}()
	close(m.queue)
	close(m.errorsQueue)
	wg.Wait()
	l.Println("[Close] - all cancelled, returning")
}

// listen starts to listen for send email requests and calling send when receive the request
func (m *Mailer) listen() {
	for {
		select {
		case <-m.done:
			l.Println("[listen] - leaving via done")
			return
		case emailReq := <-m.queue:
			go m.send(emailReq)
		}
	}
}

// listenErrors
func (m *Mailer) listenErrors() {
	for {
		select {
		case <-m.done:
			l.Println("[listenErrors] - leaving via done")
			return
		case err := <-m.errorsQueue:
			go m.callback(err)
		}
	}
}

func (m *Mailer) send(er *emailRequest) {

	l.Printf("sending email to: %s - body: %s\n", er.recipients, er.data[:50])
	time.Sleep(1 * time.Second)
	select {
	case <-m.done:
		l.Printf(" [send - done]: cancelled while sending to: %s \n", er.recipients)
		return
	// case <-time.After(time.Second * 1): // simulate sending
	default:
		l.Printf("Real mail sending to %s\n", er.recipients)
		// Sender data.
		from := os.Getenv("APP_EMAIL_ADDRESS")
		password := os.Getenv("APP_EMAIL_PASSWORD")

		// Receiver email address.
		to := er.recipients

		// smtp server configuration.
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")

		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpHost)

		l.Println("auth:")
		l.Println(auth)
		l.Println()

		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, er.data)
		if err != nil {
			m.emitError(err)
		}
	}

	if bytes.Contains(er.data, []byte("error")) {
		m.emitError(fmt.Errorf("en error ocurred while sending email to: %s", er.recipients[0]))
		return
	}
	l.Printf("email to: %s sent\n", er.recipients)
}

func (m *Mailer) emitError(err error) {
	select {
	case <-m.done:
	default:
		if err != nil {
			m.errorsQueue <- err
		}
	}

}
