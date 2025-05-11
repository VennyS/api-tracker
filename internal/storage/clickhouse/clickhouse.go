package clickhouse

import (
	"api-tracker/internal/config"
	"api-tracker/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseStorage struct {
	db *sql.DB
}

func New(cfg config.ClickHouseConfig) (*ClickHouseStorage, error) {
	// Формируем DSN в правильном формате
	hostWithPort := fmt.Sprintf("%s:%d", cfg.Addr, cfg.PortNative)

	dsn := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s",
		hostWithPort, // Теперь включает порт
		cfg.User,
		cfg.Password,
		cfg.DB,
	)

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверяем соединение с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &ClickHouseStorage{db: db}, nil
}

func (s *ClickHouseStorage) InsertLog(log models.APIRequestLog) error {
	query := `
		INSERT INTO requests (
			timestamp, 
			method, 
			path, 
			status_code, 
			latency_ms, 
			ip, 
			user_agent, 
			service_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query,
		log.Timestamp,
		log.Method,
		log.Path,
		log.StatusCode,
		log.LatencyMs,
		log.IP,
		log.UserAgent,
		log.ServiceName,
	)

	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}
