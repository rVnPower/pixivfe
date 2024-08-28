package audit

import (
	"log"

	"github.com/openzipkin/zipkin-go/model"
)

type LogReporter struct {
	l *log.Logger
}

func (rr LogReporter) Send(m model.SpanModel) {
	rr.l.Print(m.Name)
}

func (rr LogReporter) Close() error {
	return nil
}
