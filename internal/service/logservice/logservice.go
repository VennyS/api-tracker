package logservice

import (
	"api-tracker/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInternal      = errors.New("internal error")
	ErrInvalidInput  = errors.New("invalid input data")
	validHTTPMethods = map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodOptions: true,
		http.MethodHead:    true,
	}
)

type LogRepository interface {
	InsertLog(context.Context, models.APIRequestLog) error
}

type LogService struct {
	repo LogRepository
	log  *slog.Logger
}

func New(repo LogRepository, logger *slog.Logger) *LogService {
	return &LogService{
		repo: repo,
		log:  logger.With("component", "logService"),
	}
}

func (s *LogService) AddLog(ctx context.Context, log models.APIRequestLog) error {
	const op = "logService.AddLog"
	logger := s.log.With("operation", op)

	// Валидация входных данных
	if err := validateLog(log); err != nil {
		logger.Warn("validation failed",
			"error", err,
			"method", log.Method,
			"path", log.Path,
		)
		return fmt.Errorf("%s: %w: %w", op, ErrInvalidInput, err)
	}

	// Установка timestamp по умолчанию
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	start := time.Now()
	err := s.repo.InsertLog(ctx, log)
	duration := time.Since(start)

	if err != nil {
		logger.Error("failed to insert log",
			"error", err,
			"method", log.Method,
			"path", log.Path,
			"duration", duration,
		)
		return fmt.Errorf("%s: %w", op, ErrInternal)
	}

	logger.Debug("log successfully added",
		"method", log.Method,
		"path", log.Path,
		"status", log.StatusCode,
		"duration", duration,
	)

	return nil
}

func validateLog(log models.APIRequestLog) error {
	var validationErrors []string

	if log.Method == "" {
		validationErrors = append(validationErrors, "method is required")
	} else if !validHTTPMethods[strings.ToUpper(log.Method)] {
		validationErrors = append(validationErrors, fmt.Sprintf("invalid HTTP method: %s", log.Method))
	}

	if log.Path == "" {
		validationErrors = append(validationErrors, "path is required")
	}

	if log.ServiceName == "" {
		validationErrors = append(validationErrors, "service name is required")
	}

	if log.StatusCode < 100 || log.StatusCode > 599 {
		validationErrors = append(validationErrors, fmt.Sprintf("invalid HTTP status code: %d", log.StatusCode))
	}

	if log.LatencyMs < 0 {
		validationErrors = append(validationErrors, "latency must be non-negative")
	}

	if strings.TrimSpace(log.IP) == "" {
		validationErrors = append(validationErrors, "IP is required")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}
