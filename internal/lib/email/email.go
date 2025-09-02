package email

import (
	htmltemplate "html/template"
	"log"
	"os"
	texttemplate "text/template"

	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/wneessen/go-mail"
)

type SendTo struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

func Send(dates []*SendTo) error {
	textTpl, err := texttemplate.ParseFiles("internal/lib/templates/email.txt")
	if err != nil {
		return err
	}

	htmlTpl, err := htmltemplate.ParseFiles("internal/lib/templates/email.html")
	if err != nil {
		return err
	}

	var messages []*mail.Msg

	for _, date := range dates {
		message := mail.NewMsg()

		if err := message.FromFormat("System", os.Getenv("SMTP_USER")); err != nil {
			return err
		}

		if err := message.AddToFormat(date.Name, date.Email); err != nil {
			return err
		}

		message.SetDate()
		message.SetMessageID()
		message.SetBulk()
		message.Subject("Authorization data")

		if err := message.SetBodyTextTemplate(textTpl, date); err != nil {
			return err
		}

		if err := message.AddAlternativeHTMLTemplate(htmlTpl, date); err != nil {
			return err
		}

		messages = append(messages, message)
	}
	client, err := mail.NewClient(env.GetSmtpHost(),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(env.GetSmtpUser()),
		mail.WithPassword(env.GetSmtpPassword()),
	)
	if err != nil {
		return err
	}

	if err := client.DialAndSend(messages...); err != nil {
		return err
	}

	log.Println("Test mail successfully delivered.")

	return nil
}
