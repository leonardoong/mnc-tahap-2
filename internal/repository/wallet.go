package repository

import (
	"database/sql"
	"time"
)

type IWalletRepository interface {
	GetCurrentBalance(userID string) (float64, error)
	UpdateBalance(tx *sql.Tx, userID string, balance float64, updateAt time.Time) error
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) IWalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetCurrentBalance(userID string) (balance float64, err error) {
	query := `
		SELECT balance
		FROM wallets
		WHERE user_id = ?
	`
	row := r.db.QueryRow(query, userID)

	err = row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return balance, nil
		}
		return balance, err
	}

	return balance, nil
}

func (r *walletRepository) UpdateBalance(tx *sql.Tx, userID string, balance float64, updateAt time.Time) error {
	query := `
		UPDATE wallets
		SET balance = ?, updated_at = ?
		WHERE user_id = ?
	`
	_, err := tx.Exec(query, balance, updateAt, userID)

	return err
}
