package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leonardoong/e-wallet/config"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/repository"
)

type ITransactionService interface {
	StartTopUp(req *entity.PublishTopUpRequest) (string, error)
	StartPayment(req *entity.PaymentRequest) error
	StartTransfer(req *entity.TransferRequest) error

	ProcessTopUp(req entity.PublishTopUpRequest) error
	FindTopupByTopUpID(topUpID string) (*entity.Transaction, error)
}

type transactionService struct {
	config                *config.Config
	db                    *sql.DB
	transactionRepository repository.ITransactionRepository
	walletRepository      repository.IWalletRepository
}

func NewTransactionService(config *config.Config, dbConn *sql.DB, transactionRepo repository.ITransactionRepository, walletRepo repository.IWalletRepository) ITransactionService {
	return &transactionService{
		config:                config,
		db:                    dbConn,
		transactionRepository: transactionRepo,
		walletRepository:      walletRepo,
	}
}

func (s *transactionService) StartTopUp(req *entity.PublishTopUpRequest) (string, error) {
	topUpUuid := uuid.New().String()

	payload := entity.PublishTopUpRequest{
		TopUpID: topUpUuid,
		Amount:  req.Amount,
		UserID:  req.UserID,
	}

	s.transactionRepository.PublishTopUp(payload)
	return topUpUuid, nil
}

func (s *transactionService) ProcessTopUp(req entity.PublishTopUpRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// get current balance
	balanceBefore, err := s.walletRepository.GetCurrentBalance(tx, req.UserID)
	if err != nil {
		fmt.Println("error sini")
		fmt.Println(err.Error())

		return err
	}

	balanceAfter := balanceBefore + req.Amount

	now := time.Now()

	transaction := entity.Transaction{
		TransactionID: req.TopUpID,
		UserID:        req.UserID,
		Type:          "CREDIT",
		Amount:        req.Amount,
		Status:        "SUCCESS",
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.transactionRepository.InsertTopUp(tx, transaction)
	if err != nil {
		fmt.Println("error sini 2")
		fmt.Println(err.Error())
		return err
	}

	err = s.walletRepository.UpdateBalance(tx, req.UserID, balanceAfter, now)
	if err != nil {
		fmt.Println("error sini 3")
		fmt.Println(err.Error())
		return err
	}

	tx.Commit()

	return err
}

func (s *transactionService) StartPayment(req *entity.PaymentRequest) error {
	return nil
}

func (s *transactionService) StartTransfer(req *entity.TransferRequest) error {
	return nil
}

func (s *transactionService) FindTopupByTopUpID(topUpID string) (*entity.Transaction, error) {
	transaction, err := s.transactionRepository.FindTopupByTopUpID(topUpID)
	if err != nil {
		return nil, err
	}
	return transaction, err
}
