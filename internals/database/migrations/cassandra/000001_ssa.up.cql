CREATE TABLE IF NOT EXISTS ssa (
    session_id UUID PRIMARY KEY,
    start_time TIMESTAMP,
    bandwidth_usage DECIMAL,
    content_id TEXT,
    device_type TEXT,
    quality_level TEXT,
);