-- SQL migration files
CREATE TABLE IF NOT EXISTS api_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key         TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    is_active   BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key      TEXT NOT NULL,
    model        TEXT NOT NULL,
    prompt_tokens    INT,
    response_tokens  INT,
    latency_ms   INT,
    success      BOOLEAN,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);