// defines the schema for the normalised data
package intelligence

import "time"

type TelemetryEvent struct {
    ID        string         `json:"id"`
    Timestamp time.Time      `json:"timestamp"`
    Service   string         `json:"service"`
    EventType string         `json:"event_type"` // "log", "metric", "deployment", "trace"
    Severity  string         `json:"severity"`   // "info", "warn", "error", "critical"
    Message   string         `json:"message"`
    Metadata  map[string]any `json:"metadata"`
    TraceID   string         `json:"trace_id"`
}
