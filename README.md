# ObsIntel

A production-grade AI infrastructure platform built in Go. ObsIntel is a high-throughput LLM gateway and real-time intelligence system designed for AI-powered observability and incident intelligence.

> Currently in active development. Gateway layer complete. Streaming pipeline and RAG layer in progress.

---

## What it is

ObsIntel is built in layers:

- **LLM Gateway** — A reliable, observable proxy for LLM API calls with rate limiting, auth, cost tracking, semantic caching, and streaming
- **Streaming Pipeline** — Real-time event ingestion via Redpanda/Kafka for continuous data processing *(in progress)*
- **RAG Engine** — Retrieval-augmented generation over live streaming data *(in progress)*
- **Eval System** — Automated quality tracking and regression detection for AI outputs *(planned)*

The primary use case is **SRE/Incident Intelligence** — ingesting telemetry (logs, metrics, traces, deployment events) in real time and enabling engineers to query what's wrong, why, and what happened before.

---

## Gateway Features

- Multi-provider LLM support via provider interface (Gemini implemented, OpenAI-ready)
- SSE streaming responses
- Token counting and cost tracking (live pricing from OpenRouter)
- Semantic cache with pgvector — similar queries return cached responses without hitting the LLM
- Per-IP rate limiting with token bucket algorithm
- API key authentication
- Request logging with latency, tokens, cost, and error tracking
- Prometheus metrics endpoint
- Grafana dashboard

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| API Framework | Gin |
| Database | PostgreSQL + pgvector |
| Embeddings | Ollama + nomic-embed-text (local) |
| LLM Provider | Gemini API |
| Streaming | Redpanda (Kafka-compatible) |
| Observability | Prometheus + Grafana |
| Deployment | Docker |

---

