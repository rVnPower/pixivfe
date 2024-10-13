package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/soluble-ai/go-jnode"

	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

func Diagnostics(w http.ResponseWriter, r *http.Request) error {
	return RenderHTML(w, r, Data_diagnostics{})
}

func ResetDiagnosticsData(w http.ResponseWriter, r *http.Request) {
	audit.RecordedRequestSpans = audit.RecordedRequestSpans[:0]
	utils.RedirectToWhenceYouCame(w, r)
}

// formatSpanSummary creates a SpanSummary string from audit.Span
func formatSpanSummary(span audit.Span) string {
	duration := float64(audit.Duration(span)) / float64(time.Second)
	return fmt.Sprintf("%s - %s - %v - %v - %.3fs",
		span.GetStartTime().Format(time.RFC3339),
		span.Component(),
		span.Action(),
		span.Outcome(),
		duration,
	)
}

func DiagnosticsData(w http.ResponseWriter, _ *http.Request) error {
	data := jnode.NewArrayNode()
	for _, span := range audit.RecordedRequestSpans {
		bytes, err := json.Marshal(span)
		if err != nil {
			return err
		}
		obj, err := jnode.FromJSON(bytes)
		if err != nil {
			return err
		}
		obj.Put("LogLine", formatSpanSummary(span))
		data.Append(obj)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(data)
}
