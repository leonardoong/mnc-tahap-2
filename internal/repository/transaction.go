package repository

import (
	"database/sql"

	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/publisher"
)

type ITransactionRepository interface {
	InsertTopUp(tx *sql.Tx, transaction entity.Transaction) error
	Payment(transaction *entity.Transaction) error
	Transfer(transaction *entity.Transaction) error
	TransactionHistory() error

	PublishTopUp(payload entity.PublishTopUpRequest) error
	FindTopupByTopUpID(topUpID string) (*entity.Transaction, error)
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

func (r *transactionRepository) InsertTopUp(tx *sql.Tx, transaction entity.Transaction) error {
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

func (r *transactionRepository) Payment(transaction *entity.Transaction) error {
	return nil
}

func (r *transactionRepository) Transfer(transaction *entity.Transaction) error {
	return nil
}

func (r *transactionRepository) TransactionHistory() error {
	return nil
}

func (r *transactionRepository) FindTopupByTopUpID(topUpID string) (*entity.Transaction, error) {
	query := `
	SELECT *
	FROM transactions
	WHERE transaction_id = ?
	`
	row := r.db.QueryRow(query, topUpID)

	transaction := &entity.Transaction{}
	err := row.Scan(&transaction.ID, &transaction.TransactionID, &transaction.UserID, &transaction.Type,
		&transaction.Amount, &transaction.BalanceBefore, &transaction.BalanceAfter, &transaction.Description,
		&transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return transaction, nil
}
