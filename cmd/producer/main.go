package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/vaibhavvvvv/obsintel/internal/intelligence"
)

func main() {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:29092"),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// produce 5 test events
	for i := 0; i < 5; i++ {
		event := intelligence.TelemetryEvent{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Service:   "checkout-service",
			EventType: "metric",
			Severity:  "info",
			Message:   fmt.Sprintf("test event %d", i),
			Metadata: map[string]any{
				"latency_ms": 150 + i*10,
				"error_rate": 0.01,
			},
			TraceID: uuid.New().String(),
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal event: %v", err)
			continue
		}

		record := &kgo.Record{
			Topic: "telemetry.events",
			Value: data,
		}

		if err := client.ProduceSync(ctx, record).FirstErr(); err != nil {
			log.Printf("Failed to produce event: %v", err)
			continue
		}

		fmt.Printf("Produced event %d: %s\n", i, event.ID)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Done producing 5 events")
}
