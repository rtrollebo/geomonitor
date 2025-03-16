package monitor

import (
	"context"
	"errors"
	"log"
	"net/smtp"
)

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

func Run(ctx context.Context, sender string, recipient string, smtpAddress string, smtpPort, smtpPass string) error {
	logError := ctx.Value("logerror").(*log.Logger)
	logInfo := ctx.Value("loginfo").(*log.Logger)
	logInfo.Println("Sending email notification: " + recipient)
	msg := "test notification"
	err := smtp.SendMail(smtpAddress+":"+smtpPort,
		smtp.PlainAuth("", sender, smtpPass, smtpAddress),
		sender, []string{recipient}, []byte(msg))

	if err != nil {
		logError.Println("Failed to send email: " + err.Error())
		return err
	}
	return nil
}
