package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rtrollebo/geomonitor/monitor"
)

const application = "geomonitor"
var version = "undefined"

func ReadChannel(messages chan monitor.TaskResult, max_routines int) {
	timeOut := 10
	for {
		select {
		case <-messages:
		case <-time.After(time.Second * time.Duration(timeOut)):
			if runtime.NumGoroutine() > max_routines {
				fmt.Println("Routines taking too long time. Aborting.")
				os.Exit(1)
			}
			return
		}
	}
}

func RunSchedule(interval int, sign chan os.Signal, messages chan monitor.TaskResult, tasks []monitor.Runner, ctx context.Context) {

	// TODO: counter modulus interval
	for {
		select {
		case <-sign:
			fmt.Println("shutdown")
			return
		default:
			for _, v := range tasks {
				go v.Run(messages, ctx)
			}
			ReadChannel(messages, 3)
			time.Sleep(time.Duration(interval) * time.Second)

		}
	}
}

func main() {

	fileLogging, err := os.OpenFile("geomonitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fileLogging.Close()

	var LogInfo *log.Logger
	var LogWarning *log.Logger
	var LogError *log.Logger

	LogInfo = log.New(fileLogging, "geomonitor - INFO: ", log.Ldate|log.Ltime)
	LogWarning = log.New(fileLogging, "geomonitor - WARNING: ", log.Ldate|log.Ltime)
	LogError = log.New(fileLogging, "geomonitor - ERROR: ", log.Ldate|log.Ltime)

	monitorConfig, confErr := monitor.ReadConfigFile("config.json")
	if confErr != nil {
		LogError.Println("Failed to read config: " + confErr.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "loginfo", LogInfo)
	ctx = context.WithValue(ctx, "logwarning", LogWarning)
	ctx = context.WithValue(ctx, "logerror", LogError)

	LogInfo.Println("Starting application "+application+" version "+version)

	scheduleInterval := monitorConfig.TaskInterval

	task1 := monitor.TaskDefault{Url: monitorConfig.GoesServiceUrl, Name: "GoesXray"}
	task2 := monitor.TaskNotify{Sender: monitorConfig.NotifySender, Recipient: monitorConfig.NotifyRecipients[0], SmtpAddress: monitorConfig.NotifySmtpAddress, SmtpPort: monitorConfig.NotifySmtpPort, SmtpPass: monitorConfig.NotifySmtpPass, Name: "NotifyTask"}
	tasks := []monitor.Runner{task1, task2}

	os_signal := make(chan os.Signal, 1)
	signal.Notify(os_signal, syscall.SIGINT, syscall.SIGTERM)
	messages := make(chan monitor.TaskResult)

	RunSchedule(scheduleInterval, os_signal, messages, tasks, ctx)

}
