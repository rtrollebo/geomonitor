package monitor

import (
	"context"
	"time"

	"github.com/rtrollebo/geomonitor/geo"
)

type Runner interface {
	Run(comm chan TaskResult, ctx context.Context)
}

type TaskResult struct {
	TimeTaken     int64
	Completed     bool
	TimeStarted   time.Time
	TimeCompleted time.Time
	Cause         string
}

type TaskDefault struct {
	Name         string
	Url          string
	Observations []geo.GoesXray
}

type TaskNotify struct {
	Name        string
	Sender      string
	Recipient   string
	SmtpAddress string
	SmtpPort    string
	SmtpPass    string
}

func (t TaskDefault) Run(ch chan TaskResult, ctx context.Context) {
	err := geo.Run(t.Url, ctx)
	if err != nil {
		ch <- TaskResult{Cause: err.Error()}
		return
	}
	ch <- TaskResult{TimeTaken: 1, TimeStarted: time.Now(), TimeCompleted: time.Now(), Completed: false}
}

func (t TaskNotify) Run(ch chan TaskResult, ctx context.Context) {
	err := Run(ctx, t.Sender, t.Recipient, t.SmtpAddress, t.SmtpPort, t.SmtpPass)
	if err != nil {
		ch <- TaskResult{Cause: err.Error()}
		return
	}
	ch <- TaskResult{TimeTaken: 1, TimeStarted: time.Now(), TimeCompleted: time.Now(), Completed: false}
}
