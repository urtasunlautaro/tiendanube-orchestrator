package main

import (
	"log/slog"
	"net/http"
	"os"

	h "github.com/urtasunlautaro/orchestrator/internal/handler"
	p "github.com/urtasunlautaro/orchestrator/internal/processor"
	s "github.com/urtasunlautaro/orchestrator/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	dbService := s.NewDatabase(logger)
	numService := s.NewNumerator(logger)
	processor := p.NewProcessor(dbService, numService, logger)
	handler := h.NewOrchestratorHandler(processor, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /transactions", handler.CreateOperationHandler)

	logger.Info("Server starting on port 8000...")

	if err := http.ListenAndServe(":8000", mux); err != nil {
		logger.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}
