package models

import "time"

type OperationRequest struct {
	Value              string `json:"value"`
	Description        string `json:"description"`
	Method             string `json:"method"`
	CardNumber         string `json:"cardNumber"`
	CardHolderName     string `json:"cardHolderName"`
	CardExpirationDate string `json:"cardExpirationDate"`
	CardCVV            string `json:"cardCvv"`
}

type Transaction struct {
	ID                 string `json:"id"`
	Value              string `json:"value"`
	Description        string `json:"description"`
	Method             string `json:"method"`
	CardNumber         string `json:"cardNumber"`
	CardHolderName     string `json:"cardHolderName"`
	CardExpirationDate string `json:"cardExpirationDate"`
	CardCVV            string `json:"cardCvv"`
}

type Receivable struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	CreateDate    time.Time `json:"create_date"`
	PaymentDate   time.Time `json:"payment_date"`
	Subtotal      string    `json:"subtotal"`
	Discount      string    `json:"discount"`
	Total         string    `json:"total"`
}
