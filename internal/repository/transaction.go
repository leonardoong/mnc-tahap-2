package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/publisher"
)

type ITransactionRepository interface {
	InsertTransaction(tx *sql.Tx, transaction entity.Transaction) error

	PublishTopUp(payload entity.PublishTopUpRequest) error
	PublishPayment(payload entity.PaymentRequest) error
	PublishTransfer(payload entity.TransferRequest) error
	FindTransactionByID(topUpID string) (*entity.Transaction, error)
	FindTransactionsByUserID(userID string) ([]*entity.Transaction, error)
}

type transactionRepository struct {
	db             *sql.DB
	redisPublisher *publisher.Publisher
}

func NewTransactionRepository(db *sql.DB, redisPublisher *publisher.Publisher) ITransactionRepository {
	return &transactionRepository{db: db, redisPublisher: redisPublisher}
}

func (r *transactionRepository) PublishTopUp(payload entity.PublishTopUpRequest) error {
	err := r.redisPublisher.Enqueue("topup_job", work.Q{
		"top_up_id": payload.TopUpID,
		"amount":    payload.Amount,
		"user_id":   payload.UserID,
	})
	return err
}

func (r *transactionRepository) PublishPayment(payload entity.PaymentRequest) error {
	err := r.redisPublisher.Enqueue("payment_job", work.Q{
		"payment_id": payload.PaymentID,
		"amount":    payload.Amount,
		"user_id":   payload.UserID,
		"remarks" : payload.Remarks,
	})
	return err
}

func (r *transactionRepository) PublishTransfer(payload entity.TransferRequest) error {
	err := r.redisPublisher.Enqueue("transfer_job", work.Q{
		"transfer_id": payload.TransferID,
		"target_transfer_id": payload.TransferID,
		"amount":    payload.Amount,
		"user_id":   payload.UserID,
		"target_user": payload.TargetUser,
		"remarks" : payload.Remarks,
	})
	return err
}

func (r *transactionRepository) InsertTransaction(tx *sql.Tx, transaction entity.Transaction) error {
	query := `
		INSERT INTO transactions (transaction_id, user_id, type, amount, balance_before, balance_after, status, description, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.Exec(query, transaction.TransactionID, transaction.UserID, transaction.Type, transaction.Amount, transaction.BalanceBefore, transaction.BalanceAfter, transaction.Status, transaction.Description, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

func (r *transactionRepository) TransactionHistory() error {
	return nil
}

func (r *transactionRepository) FindTransactionByID(transactionID string) (*entity.Transaction, error) {
	query := `
	SELECT *
	FROM transactions
	WHERE transaction_id = ?
	`
	row := r.db.QueryRow(query, transactionID)

	transaction := &entity.Transaction{}
	var createdAtStr, updatedAtStr string
	err := row.Scan(&transaction.ID, &transaction.TransactionID, &transaction.UserID, &transaction.Type,
		&transaction.Amount, &transaction.BalanceBefore, &transaction.BalanceAfter, &transaction.Description,
		&transaction.Status, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	transaction.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	transaction.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) FindTransactionsByUserID(userID string) ([]*entity.Transaction, error) {
	query := `
	SELECT *
	FROM transactions
	WHERE user_id = ?
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.Transaction

	for rows.Next() {
		transaction := &entity.Transaction{}
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&transaction.ID, &transaction.TransactionID, &transaction.UserID, &transaction.Type,
			&transaction.Amount, &transaction.BalanceBefore, &transaction.BalanceAfter, &transaction.Description,
			&transaction.Status, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}

		transaction.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		transaction.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}