package monitor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/rtrollebo/geomonitor/geo"
	"github.com/rtrollebo/geomonitor/internal"
)

type Notifications struct {
	Time      time.Time
	Recipient string
}

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

	notifications, readFileErrorNot := internal.ReadFile[Notifications]("notifications.json")
	if readFileErrorNot != nil {
		logError.Println("Failed to read notifications file")
		return readFileErrorNot
	}

	for _, not := range notifications {
		if not.Time.After(time.Now().Add(time.Duration(-1) * time.Hour)) {
			logInfo.Println("Notification already sent")
			return nil
		}
	}

	events, readFileErr := internal.ReadFile[geo.GeoEvent]("events.json")
	if readFileErr != nil {
		logError.Println("Failed to read events file")
	}
	if events == nil || len(events) == 0 {
		logInfo.Println("No events found")
		return nil
	}

	var recentEvent geo.GeoEvent
	newNotification := false
	for _, event := range events {
		if event.Time.After(recentEvent.Time) && event.Processed && event.Time.After(time.Now().Add(time.Duration(-1)*time.Hour)) {
			recentEvent = event
			newNotification = true
		}
	}
	if !newNotification {
		logInfo.Println("No new events found")
		return nil
	}

	// Write notfications
	notifications = append(notifications, Notifications{Time: time.Now(), Recipient: recipient})
	writeErrorNot := internal.WriteFile[Notifications](notifications, "notifications.json")
	if writeErrorNot != nil {
		logError.Println("Failed to write notifications file")
		return writeErrorNot
	}

	logInfo.Println("Sending email notification: " + recipient)
	msg := "Solar flare (XRay event) occurred at " + recentEvent.Time.Format("2006-01-02 15:04:05") + " UTC\n"
	msg += fmt.Sprintf("Event: %s\n", recentEvent.Description)
	msg += fmt.Sprintf("Peak flux: %.2E\n", recentEvent.Value)
	msg += fmt.Sprintf("Category: %d\n", recentEvent.Cat)
	msg += "\n\ngeomonitor"
	err := smtp.SendMail(smtpAddress+":"+smtpPort,
		smtp.PlainAuth("", sender, smtpPass, smtpAddress),
		sender, []string{recipient}, []byte(msg))

	if err != nil {
		logError.Println("Failed to send email: " + err.Error())
		return err
	}
	return nil
}
