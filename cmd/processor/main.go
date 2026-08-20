// starts stream processor
package main

import (
    "context"
    "encoding/json"
    "log"
    "fmt"

    "github.com/twmb/franz-go/pkg/kgo"
    "github.com/vaibhavvvvv/obsintel/config"
    "github.com/vaibhavvvvv/obsintel/internal/intelligence"
    "github.com/vaibhavvvvv/obsintel/store"
    "github.com/vaibhavvvvv/obsintel/store/queries"
)

func main() {
    config.Init()
    
    ctx := context.Background()
    store.Init(ctx)
    defer store.Close()

    client, err := kgo.NewClient(
        kgo.SeedBrokers("localhost:29092"), //tells the client where Redpanda is as in a cluster there are multiple brokers
        kgo.ConsumerGroup("obsintel-processor"),
        kgo.ConsumeTopics("telemetry.events"),
    )
    if err != nil {
        log.Fatalf("Failed to create Kafka client: %v", err)
    }
    defer client.Close()

    fmt.Println("Processor started, waiting for events...")

    for {
        fetches := client.PollFetches(ctx)
        if errs := fetches.Errors(); len(errs) > 0 {
            log.Printf("Fetch errors: %v", errs)
            continue
        }

        fetches.EachRecord(func(record *kgo.Record) {
            var event intelligence.TelemetryEvent
            if err := json.Unmarshal(record.Value, &event); err != nil {
                log.Printf("Failed to unmarshal event: %v", err)
                return
            }

            if err := queries.InsertTelemetryEvent(ctx, event); err != nil {
                log.Printf("Failed to insert event %s: %v", event.ID, err)
                return
            }

            fmt.Printf("Stored event: service=%s type=%s severity=%s\n",
                event.Service, event.EventType, event.Severity)
        })
    }
}