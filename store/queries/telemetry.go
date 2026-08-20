package queries

import (
    "context"
    "encoding/json"
    "github.com/vaibhavvvvv/obsintel/internal/intelligence"
    "github.com/vaibhavvvvv/obsintel/store"
)

func InsertTelemetryEvent(ctx context.Context, event intelligence.TelemetryEvent) error {
    metadata, err := json.Marshal(event.Metadata)
    if err != nil {
        return err
    }

    _, err = store.Pool.Exec(ctx,
        `INSERT INTO telemetry_events 
        (event_id, occurred_at, service, event_type, severity, message, metadata, trace_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (event_id) DO NOTHING`,
        event.ID,
        event.Timestamp,
        event.Service,
        event.EventType,
        event.Severity,
        event.Message,
        metadata,
        event.TraceID,
    )
    return err
}