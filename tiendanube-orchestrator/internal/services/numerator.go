package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const numeratorBaseUrl = "http://localhost:3000"

type Numerator interface {
	GetTwoUniqueIDs() (string, string, error)
}

type numerator struct {
	httpClient *resty.Client
}

func NewNumerator() Numerator {
	client := resty.New().
		SetBaseURL(numeratorBaseUrl).
		SetTimeout(10 * time.Second)
	return &numerator{httpClient: client}
}

type numeratorResponse struct {
	Value int `json:"value"`
}

type numeratorReqBody struct {
	OldValue int `json:"oldValue"`
	NewValue int `json:"newValue"`
}

func (n *numerator) GetTwoUniqueIDs() (string, string, error) {
	for i := 0; i < 5; i++ {
		var currentValue numeratorResponse
		response, err := n.httpClient.R().
			SetResult(&currentValue).
			Get("/numerator")
		if err != nil {
			return "", "", fmt.Errorf("failed to get current numerator: %w", err)
		}

		if response.IsError() {
			return "", "", fmt.Errorf("http error getting numerator: %s", response.Status())
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
			return "", "", fmt.Errorf("request failed: %w", err)
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
