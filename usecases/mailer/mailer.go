package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	mime       = []byte("MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n")
	subjectKey = []byte("Subject: ")
)

type Mailer struct {
	l           *zerolog.Logger
	queue       chan *emailRequest
	done        chan struct{}
	errorsQueue chan error
	errHandler  ErrorHandler
	// sendMail is a smtp.SendMail, declated here to allow overriding in tests
	sendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type emailRequest struct {
	subject    []byte
	recipients []string
	data       []byte
}

// ErrorHandler is called when an error occurs, it guarantees that err is not nil
type ErrorHandler func(err error)

// NewMailer creates new Mailer, it setup listeners and returns created instace
func NewMailer(l *zerolog.Logger, errHandler ErrorHandler) *Mailer {
	m := &Mailer{
		l:           l,
		errHandler:  errHandler,
		queue:       make(chan *emailRequest, 5),
		errorsQueue: make(chan error, 5),
		done:        make(chan struct{}),
		sendMail:    smtp.SendMail,
	}

	go m.listenErrors()
	go m.listen()
	return m
}

// Send emits event to send an email
func (m *Mailer) Send(recipients []string, subject, data []byte) {

	select {
	case <-m.done:
		m.emitError(fmt.Errorf("mailer: already closed"))
		return
	default:
		if len(recipients) == 0 {
			m.emitError(fmt.Errorf("mailer: missing recipients"))
			return
		}
		go func() {
			select {
			case <-m.done:
				m.emitError(fmt.Errorf("mailer: already closed"))

			case m.queue <- &emailRequest{
				subject:    subject,
				recipients: recipients,
				data:       data,
			}:
				// default:
			}
		}()
	}
}

// Close stops the mailer service, it will cancell all queued tasks
// and waits for tasks in progress to complete before returning
func (m *Mailer) Close() {
	close(m.done)
	close(m.queue)
	close(m.errorsQueue)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for range m.queue {
		}
		wg.Done()
	}()
	for range m.errorsQueue {
	}

	wg.Wait()
	<-m.done
	m.l.Debug().Msg("mailer closed")
}

// listen starts to listen for send email requests and calling send when receive the request
func (m *Mailer) listen() {
	for {
		select {
		case <-m.done:
			return
		case emailReq := <-m.queue:
			go m.emitError(m.send(emailReq))
		}
	}
}

// listenErrors
func (m *Mailer) listenErrors() {
	for {
		select {
		case <-m.done:
			return
		case err := <-m.errorsQueue:
			if err != nil {
				go m.errHandler(err)
			}
		}
	}
}

func (m *Mailer) send(er *emailRequest) error {

	from := os.Getenv("APP_EMAIL_ADDRESS")
	password := os.Getenv("APP_EMAIL_PASSWORD")
	if er == nil {

		return errors.Errorf("mail request is nil")
	}
	to := er.recipients

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := er.subject
	if len(subject) == 0 {
		subject = []byte(os.Getenv("APP_NAME"))
	}

	msg := append(mime, subjectKey...)
	msg = append(msg, subject...)
	msg = append(msg, '\n')
	msg = append(msg, er.data...)

	if m.l == nil {
		return errors.Errorf("nil logget!")
	}
	m.l.Printf("mailer: about to send email to  %s", to)
	err := m.sendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	return errors.WithMessagef(err, "send email to: %s", to)
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
