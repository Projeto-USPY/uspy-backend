package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/mailjet/mailjet-apiv3-go"
)

// Mailjet is the email client object
type Mailjet struct {
	APIKey string `envconfig:"USPY_MAILJET_KEY"`
	Secret string `envconfig:"USPY_MAILJET_SECRET"`

	client *mailjet.Client
}

// Email defaults
const (
	Sender = `no-reply@uspy.me`
	Name   = `USPY`
)

// Verification
const (
	VerificationSubject     = `Verifique sua conta para usar o USPY =)`
	PasswordRecoverySubject = `Aqui está seu link de recuperação de senha do USPY =)`

	VerificationContent = `Olá! Bem vindo ao USPY!

	Por questões de segurança, precisamos que você verifique a sua conta através do seguinte link:
	
	<a href="%s">Clique aqui para verificar sua conta.</a>
	`

	PasswordRecoveryContent = `Opa =), aqui está seu link de recuperação de senha!

	Caso esse pedido não tenha sido feito por você, desconsidere esse e-mail.
	
	<a href="%s">Clique aqui para redefinir sua senha.</a>
	`
)

// Setup initializes the mailjet client using its API Key and Secret
func (m *Mailjet) Setup() {
	if m.APIKey != "" && m.Secret != "" {
		m.client = mailjet.NewMailjetClient(m.APIKey, m.Secret)
	} else {
		log.Error("failed to configure email client")
	}
}

// Send sends an email.
//
// It takes a target email recipient, subject and HTML content.
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
