package services

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	numeratorBaseUrl = "http://localhost:3000"
	retries          = 5
)

type Numerator interface {
	GetTwoUniqueIDs() (string, string, error)
}

type numerator struct {
	httpClient *resty.Client
	logger     *slog.Logger
}

func NewNumerator(logger *slog.Logger) Numerator {
	client := resty.New().
		SetBaseURL(numeratorBaseUrl).
		SetTimeout(10 * time.Second)
	return &numerator{httpClient: client, logger: logger.With("service", "numerator")}
}

type numeratorResponse struct {
	Value int `json:"value"`
}

type numeratorReqBody struct {
	OldValue int `json:"oldValue"`
	NewValue int `json:"newValue"`
}

func (n *numerator) GetTwoUniqueIDs() (string, string, error) {
	for i := 0; i < retries; i++ {
		n.logger.Info(fmt.Sprintf("getting ids, attempt #%d", i))
		var currentValue numeratorResponse
		response, err := n.httpClient.R().
			SetResult(&currentValue).
			Get("/numerator")
		if err != nil {
			msg := "failed to get current numerator"
			n.logger.Error(msg, err)
			return "", "", fmt.Errorf("%s: %w", msg, err)
		}

		if response.IsError() {
			msg := "error response getting numerator"
			n.logger.Error(msg, err)
			return "", "", fmt.Errorf("%s: %s", msg, response.Status())
		}

		oldValue := currentValue.Value
		newValue := oldValue + 2
		reqBody := numeratorReqBody{
			OldValue: oldValue,
			NewValue: newValue,
		}

		respPut, err := n.httpClient.R().
			SetBody(reqBody).
			Put("/numerator/test-and-set")
		if err != nil {
			msg := "request for numerator test-and-set failed"
			n.logger.Error(msg, err)
			return "", "", fmt.Errorf("%s: %w", msg, err)
		}

		if !respPut.IsError() {
			transactionID := strconv.Itoa(oldValue)
			receivableID := strconv.Itoa(oldValue + 1)
			return transactionID, receivableID, nil
		}

		time.Sleep(50 * time.Millisecond)
	}

	return "", "", fmt.Errorf("failed to get unique IDs")
}
