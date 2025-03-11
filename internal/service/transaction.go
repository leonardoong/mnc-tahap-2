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
	StartPayment(req *entity.PaymentRequest) (string, error)
	StartTransfer(req *entity.TransferRequest)  (*entity.StartTransferResponse, error)

	ProcessTopUp(req entity.PublishTopUpRequest) error
	ProcessPayment(req entity.PaymentRequest) error
	ProcessTransfer(req entity.TransferRequest) error

	FindTransactionByID(transactionID string) (*entity.Transaction, error)
	FindTransactionsByUserID(userID string) ([]*entity.Transaction, error)
}

type transactionService struct {
	config                *config.Config
	db                    *sql.DB
	transactionRepository repository.ITransactionRepository
	walletRepository      repository.IWalletRepository
	userRepository repository.IUserRepository
}

func NewTransactionService(config *config.Config, 
	dbConn *sql.DB, 
	transactionRepo repository.ITransactionRepository, 
	walletRepo repository.IWalletRepository,
	userRepository repository.IUserRepository) ITransactionService {
	return &transactionService{
		config:                config,
		db:                    dbConn,
		transactionRepository: transactionRepo,
		walletRepository:      walletRepo,
		userRepository: userRepository,
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
	balanceBefore, err := s.walletRepository.GetCurrentBalance(req.UserID)
	if err != nil {
		return err
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	balanceAfter := balanceBefore + req.Amount

	now := time.Now()

	topUpTransaction := entity.Transaction{
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

	err = s.transactionRepository.InsertTransaction(tx, topUpTransaction)
	if err != nil {
		return err
	}

	err = s.walletRepository.UpdateBalance(tx, req.UserID, balanceAfter, now)
	if err != nil {
		return err
	}

	tx.Commit()

	return err
}

func (s *transactionService) StartPayment(req *entity.PaymentRequest) (string, error) {
	currentBalance, err := s.walletRepository.GetCurrentBalance(req.UserID)
	if err != nil {
		return "", err
	}

	if currentBalance < req.Amount {
		return "", fmt.Errorf("Balance is not enough")
	}

	paymentUuid := uuid.New().String()
	req.PaymentID = paymentUuid

	err = s.transactionRepository.PublishPayment(*req)
	if err != nil {
		return "", err
	}
	return paymentUuid, nil
}

func (s *transactionService) ProcessPayment(req entity.PaymentRequest)(err error){
	balanceBefore, err := s.walletRepository.GetCurrentBalance(req.UserID)
	if err != nil {
		return err
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	balanceAfter := balanceBefore - req.Amount

	now := time.Now()

	paymentTransaction := entity.Transaction{
		TransactionID: req.PaymentID,
		UserID: req.UserID,
		Type: "DEBIT",
		Amount: req.Amount,
		Status: "SUCCESS",
		BalanceBefore: balanceBefore,
		BalanceAfter: balanceAfter,
		Description: req.Remarks,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.transactionRepository.InsertTransaction(tx, paymentTransaction)
	if err != nil {
		return err
	}

	err = s.walletRepository.UpdateBalance(tx, req.UserID, balanceAfter, now)
	if err != nil {
		return err
	}

	tx.Commit()

	return err
}

func (s *transactionService) StartTransfer(req *entity.TransferRequest) (*entity.StartTransferResponse, error) {
	currentBalance, err := s.walletRepository.GetCurrentBalance(req.UserID)
	if err != nil {
		return nil, err
	}

	if currentBalance < req.Amount {
		return nil, fmt.Errorf("Balance is not enough")
	}

	_, err = s.userRepository.FindByID(req.TargetUser)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Target user not found")
	} else if  err != nil {
		return nil, err
	}

	transferUuid := uuid.New().String()
	req.TransferID = transferUuid

	targetTransferUuid := uuid.New().String()
	req.TargetTransferID = targetTransferUuid

	err = s.transactionRepository.PublishTransfer(*req)
	if err != nil {
		return nil, err
	}
	return &entity.StartTransferResponse{
		TransferID: transferUuid,
		TargetTransferID: targetTransferUuid,
	}, nil
}

func (s *transactionService) ProcessTransfer(req entity.TransferRequest)(err error){
	balanceBefore, err := s.walletRepository.GetCurrentBalance(req.UserID)
	if err != nil {
		return err
	}

	targetBalanceBefore, err := s.walletRepository.GetCurrentBalance(req.TargetUser)
	if err != nil {
		return err
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	now := time.Now()
	balanceAfter := balanceBefore - req.Amount

	userTransferTransaction := entity.Transaction{
		TransactionID: req.TransferID,
		UserID: req.UserID,
		Type: "DEBIT",
		Amount: req.Amount,
		Status: "SUCCESS",
		BalanceBefore: balanceBefore,
		BalanceAfter: balanceAfter,
		Description: req.Remarks,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.transactionRepository.InsertTransaction(tx, userTransferTransaction)
	if err != nil {
		return err
	}

	err = s.walletRepository.UpdateBalance(tx, req.UserID, balanceAfter, now)
	if err != nil {
		return err
	}

	targetBalanceAfter := targetBalanceBefore + req.Amount

	targetUserTransferTransaction := entity.Transaction{
		TransactionID: req.TargetTransferID,
		UserID: req.TargetUser,
		Type: "CREDIT",
		Amount: req.Amount,
		Status: "SUCCESS",
		BalanceBefore: targetBalanceBefore,
		BalanceAfter: targetBalanceAfter,
		Description: req.Remarks,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.transactionRepository.InsertTransaction(tx, targetUserTransferTransaction)
	if err != nil {
		return err
	}

	err = s.walletRepository.UpdateBalance(tx, req.TargetUser, targetBalanceAfter, now)
	if err != nil {
		return err
	}

	tx.Commit()

	return err
}

func (s *transactionService) FindTransactionByID(topUpID string) (*entity.Transaction, error) {
	transaction, err := s.transactionRepository.FindTransactionByID(topUpID)
	if err != nil {
		return nil, err
	}
	return transaction, err
}

func (s *transactionService) FindTransactionsByUserID(userID string) ([]*entity.Transaction, error) {
	transactions, err := s.transactionRepository.FindTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return transactions, err
}
