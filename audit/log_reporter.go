package audit

import (
	"log"
	"os"
	"time"

	"github.com/openzipkin/zipkin-go/model"
)

type LogReporter struct {
	l *log.Logger
}

func (rr LogReporter) Send(m model.SpanModel) {
	rr.l.Printf("%v +%-5.3f %s", m.Timestamp.Format("2006-01-02 15:04:05.000"), float64(m.Duration)/float64(time.Second), m.Name)
}

func (rr LogReporter) Close() error {
	return nil
}

func NewLogReporter() LogReporter {
	return LogReporter{l: log.New(os.Stderr, "", 0)}
}
