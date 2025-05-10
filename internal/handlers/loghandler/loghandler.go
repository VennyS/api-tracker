package loghandler

import (
	"api-tracker/internal/domain/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

type LogService interface {
	AddLog(log models.APIRequestLog) error
}

type logHandler struct {
	logSrv LogService
}

func New(logSrv LogService) *logHandler {
	return &logHandler{logSrv: logSrv}
}

func (h *logHandler) PostLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logRequest models.APIRequestLog
		err := json.NewDecoder(r.Body).Decode(&logRequest)
		if err != nil {
			sendMessage(w, r, "Bad request", http.StatusBadRequest)
			return
		}

		err = h.logSrv.AddLog(logRequest)
		if err != nil {
			sendMessage(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		sendMessage(w, r, "ok", http.StatusOK)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request, message string, code int) {
	render.Status(r, code)
	render.JSON(w, r, map[string]string{
		"message": message,
	})
}
