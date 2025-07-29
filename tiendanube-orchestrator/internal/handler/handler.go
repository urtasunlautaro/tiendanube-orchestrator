package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/urtasunlautaro/orchestrator/internal/models"
	"github.com/urtasunlautaro/orchestrator/internal/processor"
)

type OrchestratorHandler struct {
	processor processor.Processor
	logger    *slog.Logger
}

func NewOrchestratorHandler(p processor.Processor, logger *slog.Logger) *OrchestratorHandler {
	return &OrchestratorHandler{
		processor: p,
		logger:    logger.With("handler"),
	}
}

func (h *OrchestratorHandler) CreateOperationHandler(w http.ResponseWriter, r *http.Request) {
	var request models.OperationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Warn("failed to decode request body", "error", err)
		h.writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdTransaction, err := h.processor.CreateOperation(request)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "an internal error occurred")
		return
	}

	h.writeJSON(w, http.StatusCreated, createdTransaction)
}

func (h *OrchestratorHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	errorResponse := map[string]string{"error": message}
	h.writeJSON(w, status, errorResponse)
}

func (h *OrchestratorHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to write json response", "error", err)
	}
}
