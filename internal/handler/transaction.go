package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/service"
)

type TransactionHandler struct {
	TransactionService service.ITransactionService
}

func (h *TransactionHandler) TopUp(c *gin.Context) {
	var req entity.TopUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
	}

	payload := entity.PublishTopUpRequest{
		Amount: req.Amount,
		UserID: userID.(string),
	}

	topUpID, err := h.TransactionService.StartTopUp(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id": topUpID,
		},
	})
}

func (h *TransactionHandler) FindTopUp(c *gin.Context) {
	topUpID := c.Param("top_up_id")
	if topUpID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "top up id mandatory"})
		return
	}

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
		return
	}

	transaction, err := h.TransactionService.FindTransactionByID(topUpID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id":      transaction.TransactionID,
			"amount_top_up":  transaction.Amount,
			"balance_before": transaction.BalanceBefore,
			"balance_after":  transaction.BalanceAfter,
			"created_date":   transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) Payment(c *gin.Context) {
	var req entity.PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
	}

	req.UserID = userID.(string)

	paymentID, err := h.TransactionService.StartPayment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id": paymentID,
		},
	})
}

func (h *TransactionHandler) FindPayment(c *gin.Context) {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payment id mandatory"})
		return
	}

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
		return
	}

	transaction, err := h.TransactionService.FindTransactionByID(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id":      transaction.TransactionID,
			"amount":  transaction.Amount,
			"balance_before": transaction.BalanceBefore,
			"balance_after":  transaction.BalanceAfter,
			"remarks": transaction.Description,
			"created_date":   transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req entity.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
	}

	req.UserID = userID.(string)

	transfer, err := h.TransactionService.StartTransfer(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "SUCCESS",
		"result": transfer,
	})
}

func (h *TransactionHandler) FindTransfer(c *gin.Context) {
	transferID := c.Param("transfer_id")
	if transferID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "transfer id mandatory"})
		return
	}

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
		return
	}

	transaction, err := h.TransactionService.FindTransactionByID(transferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"transfer_id":      transaction.TransactionID,
			"amount":  transaction.Amount,
			"balance_before": transaction.BalanceBefore,
			"balance_after":  transaction.BalanceAfter,
			"remarks": transaction.Description,
			"created_date":   transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) FindTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
		return
	}

	transactions, err := h.TransactionService.FindTransactionsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resultTransactions := []gin.H{}
	for _, transaction := range transactions {
		resultTransactions = append(resultTransactions,gin.H{
			"payment_id":      transaction.TransactionID,
			"amount":  	transaction.Amount,
			"balance_before": transaction.BalanceBefore,
			"balance_after":  transaction.BalanceAfter,
			"remarks": transaction.Description,
			"created_date":   transaction.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": resultTransactions,
	})
}