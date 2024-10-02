package routes

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/soluble-ai/go-jnode"

	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

func Diagnostics(w http.ResponseWriter, r *http.Request) error {
	return RenderHTML(w, r, Data_diagnostics{})
}

func ResetDiagnosticsData(w http.ResponseWriter, r *http.Request) {
	audit.RecordedSpans = audit.RecordedSpans[:0]
	utils.RedirectToWhenceYouCame(w, r)
}

func DiagnosticsData(w http.ResponseWriter, _ *http.Request) error {
	data := jnode.NewArrayNode()
	for _, span := range audit.RecordedSpans {
		bytes, err := json.Marshal(span)
		if err != nil {
			return err
		}
		obj, err := jnode.FromJSON(bytes)
		if err != nil {
			return err
		}
		obj.Put("LogLine", span.LogLine())
		data.Append(obj)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(data)
}
