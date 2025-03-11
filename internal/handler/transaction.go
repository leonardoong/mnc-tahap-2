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

	transaction, err := h.TransactionService.FindTopupByTopUpID(topUpID)
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
