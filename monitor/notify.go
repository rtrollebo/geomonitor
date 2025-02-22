package monitor

import "errors"

type Notifier interface {
	Send() error
}

type EmailNotifier struct {
	Sender    string
	Recipient string
}

type DefaultNotifier struct {
	From string
}

func (notifier EmailNotifier) Send() error {
	return errors.New("Not implemented")
}

func (notifier DefaultNotifier) Send() error {
	return errors.New("Not implemented")
}
