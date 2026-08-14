CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS semantic_cache (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key         TEXT NOT NULL,
    query_text      TEXT NOT NULL,
    query_embedding vector(768),
    response_text   TEXT NOT NULL,
    model           TEXT NOT NULL,
    hit_count       INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    last_accessed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS cache_embedding_idx 
ON semantic_cache 
USING hnsw (query_embedding vector_cosine_ops);