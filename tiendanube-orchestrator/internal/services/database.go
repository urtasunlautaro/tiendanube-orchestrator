package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urtasunlautaro/orchestrator/internal/models"
	"io"
	"log"
	"net/http"
)

const dbBaseURL = "http://localhost:8080"

func CreateTransaction(trx models.Transaction) error {
	return postJSON(dbBaseURL+"/transactions", trx)
}

func CreateReceivable(receivable models.Receivable) error {
	return postJSON(dbBaseURL+"/receivables", receivable)
}

func DeleteTransaction(id string) error {
	url := fmt.Sprintf("%s/transactions/%s", dbBaseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer closeConnection(url, resp)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to delete transaction, status: %s", resp.Status)
	}
	return nil
}

func postJSON(url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http post to %s failed: %w", url, err)
	}

	defer closeConnection(url, resp)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("received non-2xx status code from %s: %s", url, resp.Status)
	}
	return nil
}

func closeConnection(url string, resp *http.Response) {
	func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("warning: failed to close response body for url %s: %v", url, err)
		}
	}(resp.Body)
}
