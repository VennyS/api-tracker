package clickhouse

import (
	"api-tracker/internal/config"
	"api-tracker/internal/domain/models"
	"context"
	"fmt"
	"log/slog"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseStorage struct {
	db  ch.Conn
	log *slog.Logger
}

func New(cfg config.ClickHouseConfig, log *slog.Logger) (*ClickHouseStorage, error) {
	const op = "storage.clickhouse.New"
	log = log.With("op", op)

	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.PortNative)

	conn, err := ch.Open(&ch.Options{
		Addr: []string{addr},
		Auth: ch.Auth{
			Database: cfg.DB,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		Settings: map[string]interface{}{
			"send_logs_level": "trace",
		},
	})
	if err != nil {
		log.Error("connection failed", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		log.Error("ping failed", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully connected to clickhouse",
		"addr", addr,
		"db", cfg.DB,
	)

	return &ClickHouseStorage{db: conn, log: log}, nil
}

func (s *ClickHouseStorage) InsertLog(ctx context.Context, log models.APIRequestLog) error {
	const op = "storage.clickhouse.InsertLog"
	logger := s.log.With("op", op)

	query := `
		INSERT INTO api_request_logs (
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

	start := time.Now()
	if err := s.db.Exec(ctx, query,
		log.Timestamp,
		log.Method,
		log.Path,
		log.StatusCode,
		log.LatencyMs,
		log.IP,
		log.UserAgent,
		log.ServiceName,
	); err != nil {
		logger.Error("insert failed",
			"error", err,
			"method", log.Method,
			"path", log.Path,
			"duration", time.Since(start),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("log inserted",
		"method", log.Method,
		"path", log.Path,
		"status", log.StatusCode,
		"duration", time.Since(start),
	)

	return nil
}

func (s *ClickHouseStorage) Close() error {
	const op = "storage.clickhouse.Close"
	s.log.With("op", op).Info("closing clickhouse connection")
	return s.db.Close()
}
