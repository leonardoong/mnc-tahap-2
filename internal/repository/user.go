package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/leonardoong/e-wallet/internal/domain/entity"
)

type IUserRepository interface {
	Register(user *entity.User) error
	FindByPhoneNumber(phoneNumber string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	Update(user entity.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Register(user *entity.User) error {
	query := `
		INSERT INTO users (user_id, first_name, last_name, phone_number, pin, address, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, user.UserID, user.FirstName, user.LastName, user.PhoneNumber, user.Pin, user.Address, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	walletQuery := `
		INSERT INTO wallets (user_id, balance, created_at, updated_at)
		VALUES (?, 0.00, ?, ?)
	`
	_, err = r.db.Exec(walletQuery, user.UserID, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) FindByPhoneNumber(phoneNumber string) (user *entity.User, err error) {
	query := `
		SELECT *
		FROM users 
		WHERE phone_number = ?
	`
	row := r.db.QueryRow(query, phoneNumber)

	user = &entity.User{}
	var createdAtStr, updatedAtStr string
	err = row.Scan(&user.ID, &user.UserID, &user.PhoneNumber, &user.Pin, &user.FirstName, &user.LastName, &user.Address, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByID(id string) (user *entity.User, err error) {
	query := `
		SELECT *
		FROM users 
		WHERE user_id = ?
	`
	row := r.db.QueryRow(query, id)

	user = &entity.User{}
	var createdAtStr, updatedAtStr string
	err = row.Scan(&user.ID, &user.UserID, &user.PhoneNumber, &user.Pin, &user.FirstName, &user.LastName, &user.Address, &createdAtStr, &updatedAtStr)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(user entity.User) error {
	query := `
		UPDATE users 
		SET first_name = ?, last_name = ?, address = ?, updated_at = ?
		WHERE user_id = ?
	`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.Address, user.UpdatedAt, user.UserID)

	return err
}
