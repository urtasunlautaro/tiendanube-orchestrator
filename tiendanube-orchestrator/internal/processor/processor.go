package processor

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/urtasunlautaro/orchestrator/internal/models"
	"github.com/urtasunlautaro/orchestrator/internal/services"
)

var fees = map[string]float64{
	"debit_card":  0.02,
	"credit_card": 0.04,
}

type Processor interface {
	CreateOperation(input models.OperationRequest) (*models.Transaction, error)
}

type processor struct {
	dbService  services.Database
	numService services.Numerator
	logger     *slog.Logger
}

func NewProcessor(db services.Database, num services.Numerator, logger *slog.Logger) Processor {
	return &processor{
		dbService:  db,
		numService: num,
		logger:     logger.With("processor"),
	}
}

func (p *processor) CreateOperation(input models.OperationRequest) (*models.Transaction, error) {
	trxId, receivableId, err := p.numService.GetTwoUniqueIDs()
	if err != nil {
		return nil, fmt.Errorf("could not get unique IDs: %w", err)
	}

	transaction := getTrx(input, trxId)

	receivable := getReceivable(transaction, receivableId)

	if err := p.dbService.CreateTransaction(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction, aborting: %w", err)
	}

	if err := p.dbService.CreateReceivable(receivable); err != nil {
		p.logger.Warn(fmt.Sprintf("failed to create receivable, rolling back trx with id %s", trxId))
		if rollbackErr := p.dbService.DeleteTransaction(trxId); rollbackErr != nil {
			p.logger.Error(fmt.Sprintf("failed to rollback trx with id %s: %v", trxId, rollbackErr))
		}
		return nil, fmt.Errorf("failed to create receivable: %w", err)
	}

	return &transaction, nil
}

func getTrx(input models.OperationRequest, trxId string) models.Transaction {
	return models.Transaction{
		ID:                 trxId,
		Value:              input.Value,
		Description:        input.Description,
		Method:             input.Method,
		CardNumber:         input.CardNumber[len(input.CardNumber)-4:],
		CardHolderName:     input.CardHolderName,
		CardExpirationDate: input.CardExpirationDate,
		CardCVV:            input.CardCVV,
	}
}

func getReceivable(transaction models.Transaction, receivableId string) models.Receivable {
	valueFloat, _ := strconv.ParseFloat(transaction.Value, 64)
	fee := valueFloat * fees[transaction.Method]
	total := valueFloat - fee

	receivable := models.Receivable{
		ID:            receivableId,
		TransactionID: transaction.ID,
		Subtotal:      transaction.Value,
		Discount:      fmt.Sprintf("%.2f", fee),
		Total:         fmt.Sprintf("%.2f", total),
		CreateDate:    time.Now(),
	}

	if transaction.Method == "debit_card" {
		receivable.Status = "paid"
		receivable.PaymentDate = time.Now()
	} else {
		receivable.Status = "waiting_funds"
		receivable.PaymentDate = time.Now().AddDate(0, 0, 30)
	}

	return receivable
}
