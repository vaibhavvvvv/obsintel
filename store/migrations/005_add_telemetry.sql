CREATE TABLE IF NOT EXISTS telemetry_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        TEXT UNIQUE NOT NULL,
    occurred_at     TIMESTAMPTZ NOT NULL,
    service         TEXT NOT NULL,
    event_type      TEXT NOT NULL,
    severity        TEXT NOT NULL,
    message         TEXT NOT NULL,
    metadata        JSONB DEFAULT '{}',
    trace_id        TEXT,
    processed_at    TIMESTAMPTZ DEFAULT NOW()
);

-- most queries filter by service + time
CREATE INDEX idx_telemetry_service_time
ON telemetry_events(service, occurred_at DESC);

-- anomaly detection queries filter by severity + time
CREATE INDEX idx_telemetry_severity_time
ON telemetry_events(severity, occurred_at DESC);

-- trace correlation
CREATE INDEX idx_telemetry_trace
ON telemetry_events(trace_id)
WHERE trace_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS anomalies (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    detected_at         TIMESTAMPTZ DEFAULT NOW(),
    service             TEXT NOT NULL,
    anomaly_type        TEXT NOT NULL,
    severity            TEXT NOT NULL,
    description         TEXT NOT NULL,
    evidence            JSONB DEFAULT '{}',
    correlated_events   UUID[],
    deployment_correlated BOOLEAN DEFAULT false,
    resolved            BOOLEAN DEFAULT false,
    resolved_at         TIMESTAMPTZ
);

CREATE INDEX idx_anomalies_service_time
ON anomalies(service, detected_at DESC);

CREATE INDEX idx_anomalies_unresolved
ON anomalies(resolved, detected_at DESC)
WHERE resolved = false;