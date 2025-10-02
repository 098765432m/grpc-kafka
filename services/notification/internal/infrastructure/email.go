package notification_infrastructure

import "github.com/resend/resend-go/v2"

type EmailSender struct {
	client *resend.Client
}

func NewEmailSender(apiKey string) *EmailSender {
	return &EmailSender{
		client: resend.NewClient(apiKey),
	}
}

func (es *EmailSender) SendEmail(to string, subject string, body string) error {
	_, err := es.client.Emails.Send(
		&resend.SendEmailRequest{
			From:    "sadasd@com.com",
			To:      []string{to},
			Subject: subject,
			Html:    body,
		},
	)

	return err
}
