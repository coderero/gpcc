CREATE TABLE IF NOT EXISTS ssa (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    bandwidth_usage DECIMAL(10, 2) NOT NULL,
    content_id VARCHAR(255) NOT NULL,
    device_type VARCHAR(255) NOT NULL,
    quality_level VARCHAR(255) NOT NULL
);