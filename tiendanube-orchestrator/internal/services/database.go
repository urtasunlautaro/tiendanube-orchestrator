package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/urtasunlautaro/orchestrator/internal/models"
)

const dbBaseURL = "http://localhost:8080"

type Database interface {
	CreateTransaction(trx models.Transaction) error
	CreateReceivable(receivable models.Receivable) error
	DeleteTransaction(id string) error
}

type database struct {
	httpClient *resty.Client
	logger     *slog.Logger
}

func NewDatabase(logger *slog.Logger) Database {
	client := resty.New().
		SetBaseURL(dbBaseURL).
		SetTimeout(15 * time.Second)

	return &database{httpClient: client, logger: logger.With("service", "database")}
}

func (d *database) CreateTransaction(trx models.Transaction) error {
	d.logger.Info("creating trx")
	resp, err := d.httpClient.R().
		SetBody(trx).
		Post("/transactions")
	if err != nil {
		msg := "error creating trx"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s: %w", msg, err)
	}

	if resp.IsError() {
		msg := "error response creating trx"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s: %s", msg, resp.Status())
	}

	d.logger.Info("trx creation successful")
	return nil
}

func (d *database) CreateReceivable(receivable models.Receivable) error {
	d.logger.Info("creating receivable")
	resp, err := d.httpClient.R().
		SetBody(receivable).
		Post("/receivables")
	if err != nil {
		msg := "error creating receivable"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s: %w", msg, err)
	}

	if resp.IsError() {
		msg := "error response creating receivable"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s: %s", msg, resp.Status())
	}

	d.logger.Info("receivable creation successful")
	return nil
}

func (d *database) DeleteTransaction(id string) error {
	d.logger.Info(fmt.Sprintf("deleting trx with id %s", id))

	resp, err := d.httpClient.R().
		Delete(fmt.Sprintf("/transactions/%s", id))
	if err != nil {
		msg := "error deleting transaction"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s: %w", msg, err)
	}

	if resp.IsError() {
		msg := "error response deleting trx"
		d.logger.Error(msg, err)
		return fmt.Errorf("%s, status: %s", msg, resp.Status())
	}

	d.logger.Info("trx deleted")
	return nil
}
