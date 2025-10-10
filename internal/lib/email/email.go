package email

import (
	"embed"
	htmltemplate "html/template"
	"log"
	texttemplate "text/template"

	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/wneessen/go-mail"
)

type SendTo struct {
	Name     string
	Email    string
	Username string
	Password string
}

//go:embed templates/*.txt templates/*.html
var templatesFS embed.FS

func Send(data []*SendTo) error {
	textFS, err := texttemplate.ParseFS(templatesFS, "templates/email.txt")
	if err != nil {
		return err
	}
	textTpl := texttemplate.Must(textFS, err)

	htmlFS, err := htmltemplate.ParseFS(templatesFS, "templates/email.html")
	if err != nil {
		return err
	}
	htmlTpl := htmltemplate.Must(htmlFS, err)

	var messages []*mail.Msg

	for _, d := range data {
		if d.Name == "" {
			d.Name = d.Username
		}

		message := mail.NewMsg()

		if err := message.FromFormat("System", env.GetSmtpUser()); err != nil {
			return err
		}

		if err := message.AddToFormat(d.Name, d.Email); err != nil {
			return err
		}

		message.SetDate()
		message.SetMessageID()
		message.SetBulk()
		message.Subject("Authorization data")

		if err := message.SetBodyTextTemplate(textTpl, d); err != nil {
			return err
		}

		if err := message.AddAlternativeHTMLTemplate(htmlTpl, d); err != nil {
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
