package sentrylog

import (
	"fmt"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

type Logging struct {
	Log *log.Logger
}

func New(logger *log.Logger) *Logging {
	return new(Logging).init(logger)
}

func (l *Logging) init(logger *log.Logger) *Logging {
	l.Log = logger
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	return l
}

// Info log info
func (l *Logging) Info(req ...interface{}) {
	fmt.Println(req...)
}

// Error log info
func (l *Logging) Error(req ...interface{}) {
	err := fmt.Errorf(fmt.Sprintf("%v", req))
	raven.CaptureErrorAndWait(err, nil)
	fmt.Println(req...)
	return
}

func (l *Logging) ErrorWithExtra(err error, extra map[string]interface{}) {
	raven.CaptureError(raven.WrapWithExtra(err, extra), nil)
	fmt.Println(extra)
	return
}
