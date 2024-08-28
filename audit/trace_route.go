package audit

import (
	"log"
	"time"
)

type RoutePerf struct {
	StartTime   time.Time
	EndTime     time.Time
	RemoteAddr  string
	Method      string
	Path        string
	Status      int
	Err         error
	SkipLogging bool
}

func TraceRoute(data RoutePerf) {
	if data.Err != nil {
		log.Printf("Internal Server Error: %s", data.Err)
	}

	if !data.SkipLogging {
		// todo: log.Printf("%v +%v %v %v %v %v %v", time, latency, ip, method, path, status, err)
	}
}
