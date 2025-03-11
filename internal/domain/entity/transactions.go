package entity

import "time"

type Transaction struct {
	ID            uint      `json:"id"`
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount"`
}

type TopUpResponse struct {
	TopUpID       string    `json:"top_up_id"`
	AmountTopUp   float64   `json:"amount_top_up"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}

type PublishTopUpRequest struct {
	TopUpID string  `json:"top_up_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
}

type PaymentRequest struct {
	PaymentID string  `json:"payment_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
	Remarks string  `json:"remarks"`
}

type PaymentResponse struct {
	PaymentID     string    `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferRequest struct {
	TransferID string `json:"transfer_id"`
	TargetTransferID string `json:"target_transfer_id"`
	UserID  string  `json:"user_id"`
	TargetUser string  `json:"target_user"`
	Amount     float64 `json:"amount"`
	Remarks    string  `json:"remarks"`
}

type StartTransferResponse struct {
	TransferID string `json:"transfer_id"`
	TargetTransferID string `json:"target_transfer_id"`
}

type TransferResponse struct {
	TransferID    string    `json:"transfer_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
