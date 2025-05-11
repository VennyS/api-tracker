CREATE TABLE IF NOT EXISTS api_request_logs (
    timestamp     DateTime,
    method        String,
    path          String,
    status_code   UInt16,
    latency_ms    Float64,
    ip            String,
    user_agent    String,
    service_name  String
) 
ENGINE = MergeTree()
ORDER BY (timestamp, service_name);