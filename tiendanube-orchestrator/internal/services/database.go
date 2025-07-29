package services

import (
	"fmt"
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
}

func NewDatabase() Database {
	client := resty.New().
		SetBaseURL(dbBaseURL).
		SetTimeout(15 * time.Second)

	return &database{httpClient: client}
}

// Ahora los métodos usan el cliente resty del struct (d.httpClient)
func (d *database) CreateTransaction(trx models.Transaction) error {
	resp, err := d.httpClient.R().
		SetBody(trx).
		Post("/transactions") // Usamos path relativo porque ya configuramos BaseURL

	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create transaction, status: %s", resp.Status())
	}
	return nil
}

func (d *database) CreateReceivable(receivable models.Receivable) error {
	resp, err := d.httpClient.R().
		SetBody(receivable).
		Post("/receivables")

	if err != nil {
		return fmt.Errorf("error creating receivable: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create receivable, status: %s", resp.Status())
	}
	return nil
}

func (d *database) DeleteTransaction(id string) error {
	// La ruta se construye usando el path y el ID.
	// Resty se encarga de unirlo con la BaseURL.
	resp, err := d.httpClient.R().
		Delete(fmt.Sprintf("/transactions/%s", id))

	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to delete transaction, status: %s", resp.Status())
	}
	return nil
}
