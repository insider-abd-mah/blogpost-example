package log

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"time"
)

var ErrChan = make(chan error, 100)

// Initialize sentry logging
func Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		AttachStacktrace: true,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			select {
			case e := <-ErrChan:
				if e != nil {
					Error(e)
				}
			}
		}
	}()
}

// Error
func Error(err error) {
	if os.Getenv("GIN_MODE") == "debug" {
		fmt.Println(err)

		return
	}

	sentry.CaptureException(err)
}

// Errors
func Errors(errs []error) {
	for i, err := range errs {
		if i < 10 {
			ErrChan <- err
		}
	}
}

// Fatal errors when app should halt
func Fatal(err error) {
	if os.Getenv("GIN_MODE") == "debug" {
		panic(err)
	}

	sentry.Logger.Panic(err)
}
