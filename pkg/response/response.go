package response

import (
	"net/http"

	"mini-payment-system/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Count int `json:"count,omitempty" example:"2"`
}

type SuccessBody struct {
	Success bool  `json:"success" example:"true"`
	Data    any   `json:"data,omitempty"`
	Meta    *Meta `json:"meta,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code" example:"INVALID_INPUT"`
	Message string `json:"message" example:"name is required"`
}

type ErrorBody struct {
	Success bool        `json:"success" example:"false"`
	Error   ErrorDetail `json:"error"`
}

type HealthData struct {
	Status string `json:"status" example:"ok"`
}

type MessageData struct {
	Message string `json:"message" example:"resource deleted successfully"`
}

type HealthResponse struct {
	Success bool       `json:"success" example:"true"`
	Data    HealthData `json:"data"`
}

type MessageResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    MessageData `json:"data"`
}

type UserResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    entities.User `json:"data"`
}

type UserListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []entities.User `json:"data"`
	Meta    Meta            `json:"meta"`
}

type AccountResponse struct {
	Success bool             `json:"success" example:"true"`
	Data    entities.Account `json:"data"`
}

type AccountListResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    []entities.Account `json:"data"`
	Meta    Meta               `json:"meta"`
}

type TransactionResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    entities.Transaction `json:"data"`
}

type TransactionListResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    []entities.Transaction `json:"data"`
	Meta    Meta                   `json:"meta"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessBody{Success: true, Data: data})
}

func List(c *gin.Context, status int, data any, count int) {
	c.JSON(status, SuccessBody{
		Success: true,
		Data:    data,
		Meta:    &Meta{Count: count},
	})
}

func Message(c *gin.Context, status int, message string) {
	c.JSON(status, SuccessBody{
		Success: true,
		Data: MessageData{
			Message: message,
		},
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorBody{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "INVALID_REQUEST", message)
}
