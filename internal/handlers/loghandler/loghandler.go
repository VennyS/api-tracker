package loghandler

import (
	"api-tracker/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type LogService interface {
	AddLog(context.Context, models.APIRequestLog) error
}

type logHandler struct {
	log    *slog.Logger
	logSrv LogService
}

func New(logSrv LogService, logger *slog.Logger) *logHandler {
	return &logHandler{
		logSrv: logSrv,
		log:    logger.With("handler", "logHandler"),
	}
}

func (h *logHandler) PostLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logRequest models.APIRequestLog

		if err := json.NewDecoder(r.Body).Decode(&logRequest); err != nil {
			if errors.Is(err, io.EOF) {
				h.log.Warn("empty request body", "method", r.Method, "path", r.URL.Path)
			} else {
				h.log.Error("failed to decode request body", "error", err, "method", r.Method, "path", r.URL.Path)
			}
			sendMessage(w, r, "Bad request", http.StatusBadRequest)
			return
		}

		if err := h.logSrv.AddLog(r.Context(), logRequest); err != nil {
			h.log.Error("failed to insert log",
				"error", err,
				"method", logRequest.Method,
				"path", logRequest.Path,
			)
			sendMessage(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		h.log.Debug("log successfully added",
			"method", logRequest.Method,
			"path", logRequest.Path,
			"status_code", logRequest.StatusCode,
		)
		sendMessage(w, r, "ok", http.StatusOK)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request, message string, code int) {
	render.Status(r, code)
	render.JSON(w, r, map[string]string{
		"message": message,
	})
}
