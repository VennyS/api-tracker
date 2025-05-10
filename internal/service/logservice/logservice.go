package logservice

import (
	"api-tracker/internal/domain/models"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInternal = errors.New("internal error")

type LogRepository interface {
	InsertLog(log models.APIRequestLog) error
}

type logService struct {
	logRepo LogRepository
}

func New(logRepo LogRepository) *logService {
	return &logService{logRepo: logRepo}
}

func (srv *logService) AddLog(log models.APIRequestLog) error {
	if err := validateLog(log); err != nil {
		return err
	}

	if err := srv.logRepo.InsertLog(log); err != nil {
		return ErrInternal
	}

	return nil
}

func validateLog(log models.APIRequestLog) error {
	// Обязательные поля
	if log.Method == "" {
		return errors.New("method is required")
	}

	if log.Path == "" {
		return errors.New("path is required")
	}

	if log.ServiceName == "" {
		return errors.New("service name is required")
	}

	// Валидный HTTP-метод
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"PATCH": true, "OPTIONS": true, "HEAD": true,
	}
	if !validMethods[strings.ToUpper(log.Method)] {
		return fmt.Errorf("invalid HTTP method: %s", log.Method)
	}

	// Статус-код в пределах допустимого HTTP
	if log.StatusCode < 100 || log.StatusCode > 599 {
		return fmt.Errorf("invalid HTTP status code: %d", log.StatusCode)
	}

	// Латентность не может быть отрицательной
	if log.LatencyMs < 0 {
		return errors.New("latency must be non-negative")
	}

	// IP можно просто проверить на пустоту (глубокая валидация опциональна)
	if strings.TrimSpace(log.IP) == "" {
		return errors.New("IP is required")
	}

	// Можно добавить лёгкую проверку формата IP (если хочешь, скажу как через net.ParseIP)

	// Таймстемп — если нулевой, можно выставить текущий
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	return nil
}
