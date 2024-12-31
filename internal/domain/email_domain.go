package domain

type Email interface {
	Send(to, subject, body string) error
}
