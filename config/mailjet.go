package config

import (
	"github.com/mailjet/mailjet-apiv3-go"
)

type Mailjet struct {
	APIKey string `envconfig:"USPY_MAILJET_KEY" required:"true"`
	Secret string `envconfig:"USPY_MAILJET_SECRET" required:"true"`

	client *mailjet.Client
}

// Email defaults
const (
	Sender = `no-reply@uspy.me`
	Name   = `USPY`
)

// Verification
const (
	VerificationSubject = `Verifique sua conta para usar o USPY =)`
	VerificationContent = `Olá! Bem vindo ao USPY!

	Por questões de segurança, precisamos que você verifique a sua conta através do seguinte link:
	
	<a href="%s">Clique aqui para verificar sua conta.</a>
	`
)

func (m *Mailjet) Setup() {
	m.client = mailjet.NewMailjetClient(m.APIKey, m.Secret)
}

func (m *Mailjet) Send(target, subject, content string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: Sender,
				Name:  Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: target,
				},
			},
			Subject:  subject,
			HTMLPart: content,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.client.SendMailV31(&messages)
	return err
}