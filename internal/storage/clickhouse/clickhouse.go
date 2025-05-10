package clickhouse

import (
	"api-tracker/internal/config"
	"api-tracker/internal/domain/models"
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseStorage struct {
	db *sql.DB
}

func NewClickHouseStorage(cfg config.ClickHouseConfig) (*ClickHouseStorage, error) {
	dsn := fmt.Sprintf("clickhouse://%s?database=%s", cfg.Addr, cfg.DB)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &ClickHouseStorage{db: db}, nil
}

func (s *ClickHouseStorage) InsertLog(log models.APIRequestLog) error {
	query := `
		INSERT INTO requests (timestamp, method, path, status_code, latency_ms, ip, user_agent, service_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query,
		log.Timestamp,
		log.Method,
		log.Path,
		log.StatusCode,
		log.LatencyMs,
		log.IP,
		log.UserAgent,
		log.ServiceName,
	)
	return err
}
