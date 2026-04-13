package handlers

import (
	"net/http"

	"mini-payment-system/internal/delivery/http/dto"
	transactionuc "mini-payment-system/internal/usecase/transaction"
	"mini-payment-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *transactionuc.Service
}

func NewTransactionHandler(service *transactionuc.Service) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Create godoc
// @Summary Transfer money between accounts
// @Description Create a money transfer between two accounts using a single database transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transfer payload"
// @Success 201 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 409 {object} response.ErrorBody
// @Failure 422 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	transaction, err := h.service.CreateTransfer(c.Request.Context(), req.FromAccountID, req.ToAccountID, req.Amount, req.Reference)
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusCreated, transaction)
}

// List godoc
// @Summary List transactions
// @Description Return all transaction records ordered by latest created first.
// @Tags transactions
// @Produce json
// @Success 200 {object} response.TransactionListResponse
// @Failure 500 {object} response.ErrorBody
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	transactions, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.List(c, http.StatusOK, transactions, len(transactions))
}

// GetByID godoc
// @Summary Get transaction by ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.TransactionResponse
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	transaction, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusOK, transaction)
}
