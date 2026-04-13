package handlers

import (
	"net/http"

	"mini-payment-system/internal/delivery/http/dto"
	accountuc "mini-payment-system/internal/usecase/account"
	"mini-payment-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service *accountuc.Service
}

func NewAccountHandler(service *accountuc.Service) *AccountHandler {
	return &AccountHandler{service: service}
}

// Create godoc
// @Summary Create account
// @Description Create a new account for an existing user.
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountRequest true "Create account payload"
// @Success 201 {object} response.AccountResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 409 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	account, err := h.service.Create(c.Request.Context(), req.UserID, req.InitialBalance, req.Currency)
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusCreated, account)
}

// List godoc
// @Summary List accounts
// @Description Return all accounts ordered by latest created first.
// @Tags accounts
// @Produce json
// @Success 200 {object} response.AccountListResponse
// @Failure 500 {object} response.ErrorBody
// @Router /accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	accounts, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.List(c, http.StatusOK, accounts, len(accounts))
}

// GetByID godoc
// @Summary Get account by ID
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.AccountResponse
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetByID(c *gin.Context) {
	account, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusOK, account)
}

// Update godoc
// @Summary Update account
// @Description Update account data such as currency.
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body dto.UpdateAccountRequest true "Update account payload"
// @Success 200 {object} response.AccountResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /accounts/{id} [put]
func (h *AccountHandler) Update(c *gin.Context) {
	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	account, err := h.service.Update(c.Request.Context(), c.Param("id"), req.Currency)
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusOK, account)
}

// Delete godoc
// @Summary Delete account
// @Description Delete an account by ID.
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.MessageResponse
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.Message(c, http.StatusOK, "account deleted successfully")
}
