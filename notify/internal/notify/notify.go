package notify

import (
	"bytes"
	"context"
	"log"
	"os"
	"text/template"

	mailer "github.com/shanehowearth/bcg/notify/integration/email"
	grpcProto "github.com/shanehowearth/bcg/notify/integration/grpc/proto/v1"
	messenger "github.com/shanehowearth/bcg/notify/integration/sms"
)

// Server -
type Server struct {
	SMS  messenger.Messenger
	Mail mailer.Mailer
}

// NewNotifyService -
func NewNotifyService(sms messenger.Messenger, mail mailer.Mailer) *Server {
	if sms == nil {
		log.Panic("NewNotifyService has no sms service to use")
	}
	if mail == nil {
		log.Panic("NewNotifyService has no mail service to use")
	}
	s := Server{SMS: sms, Mail: mail}
	return &s
}

// SendMail -
func (s *Server) SendMail(ctx context.Context, det *grpcProto.CustomerDetails) error {
	sender, found := os.LookupEnv("EmailSender")
	if !found {
		sender = "varun.verma@bcgdv.com"
	}

	const email = `
Dear {{.Name}},

Thank you for adding your details to our database.
{{with .Address -}}
We have your address as:
{{.}}
{{end}}
{{with .Phone -}}
We have your phone number as {{.}}.
{{end}}


Regards,
Varun Verma

	`

	t := template.Must(template.New("Welcome").Parse(email))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, det); err != nil {
		return err
	}

	// Send email
	return s.Mail.Send(det.Email, tpl.String(), sender)
}

// SendSMS -
func (s *Server) SendSMS(ctx context.Context, det *grpcProto.CustomerDetails) error {
	// TODO - Get a message from the endpoint
	message := "Congratulations on winning the top prize!"
	// Send SMS
	return s.SMS.Send(det, message)
}
