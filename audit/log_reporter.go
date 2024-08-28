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
	rr.l.Printf("%s | %v +%.3fs", m.Name,  m.Timestamp.Format(time.RFC3339), float64(m.Duration) / float64(time.Second) )
}

func (rr LogReporter) Close() error {
	return nil
}

func NewLogReporter() LogReporter {
	return LogReporter{l: log.New(os.Stderr, "", 0)}
}