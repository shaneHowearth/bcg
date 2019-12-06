package mailer

type Mailer interface {
	Send(recipient, message, sender string) error
}
