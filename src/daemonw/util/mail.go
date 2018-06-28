package util

import (
	"gopkg.in/gomail.v2"
	"errors"
)

type MailMessage struct {
	*gomail.Message
	*gomail.Dialer
}

type MailAuth struct {
	SMTPAddr string
	Port     int
	Username string
	Password string
}

type MailInfo struct {
	Subject     string
	From        string
	To          string
	Cc          map[string]string
	Body        string
	ContentType string
	Attaches    []string
}

func NewMailMessage() *MailMessage {
	m := gomail.NewMessage()
	return &MailMessage{Message: m}
}

func (m *MailMessage) SetMessageInfo(info *MailInfo) {
	m.SetSender(info.From)
	m.SetRecipient(info.To)
	if info.Cc != nil {
		for addr, name := range info.Cc {
			m.AddCc(addr, name)
		}
	}
	m.SetSubject(info.Subject)
	m.SetBody(info.ContentType, info.Body)
	if info.Attaches != nil {
		for _, attachment := range info.Attaches {
			m.AddAttach(attachment)
		}
	}
}

func (m *MailMessage) SetAuth(auth MailAuth) {
	m.Dialer = gomail.NewDialer(auth.SMTPAddr, auth.Port, auth.Username, auth.Password)
}

func (m *MailMessage) SetSender(sender string) {
	m.Message.SetHeader("From", sender)
}

func (m *MailMessage) SetRecipient(receiver string) {
	m.Message.SetHeader("To", receiver)
}

func (m *MailMessage) AddCc(address, name string) {
	recipients := m.Message.GetHeader("Cc")
	r := m.Message.FormatAddress(address, name)
	recipients = append(recipients, r)
	m.Message.SetHeader("Cc", recipients...)
}

func (m *MailMessage) SetSubject(subject string) {
	m.Message.SetHeader("Subject", subject)
}

func (m *MailMessage) SetBody(contentType, body string) {
	m.Message.SetBody(contentType, body)
}

func (m *MailMessage) AddAttach(filePath string) {
	m.Message.Attach(filePath)
}

func (m *MailMessage) Send() error {

	if m.Dialer == nil {
		return errors.New("can not dial to smtp server")
	}
	return m.Dialer.DialAndSend(m.Message)
}
