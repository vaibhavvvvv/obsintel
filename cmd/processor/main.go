// starts stream processor
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/twmb/franz-go/pkg/kgo"
    "github.com/vaibhavvvvv/obsintel/internal/intelligence"
)

func main() {
    client, err := kgo.NewClient(
        kgo.SeedBrokers("localhost:29092"), //tells the client where Redpanda is as in a cluster there are multiple brokers
        kgo.ConsumerGroup("obsintel-processor"), 
        kgo.ConsumeTopics("telemetry.events"),
    )
    if err != nil {
        log.Fatalf("Failed to create Kafka client: %v", err)
    }
    defer client.Close()

    ctx := context.Background()
    fmt.Println("Consumer started, waiting for events...")

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
            fmt.Printf("Received event: service=%s type=%s severity=%s message=%s\n",
                event.Service, event.EventType, event.Severity, event.Message)
        })
    }
}