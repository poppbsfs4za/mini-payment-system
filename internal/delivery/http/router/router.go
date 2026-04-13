package router

import (
	"mini-payment-system/internal/delivery/http/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	User        *handlers.UserHandler
	Account     *handlers.AccountHandler
	Transaction *handlers.TransactionHandler
}

func New(h Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", h.User.Create)
		v1.GET("/users", h.User.List)
		v1.GET("/users/:id", h.User.GetByID)
		v1.PUT("/users/:id", h.User.Update)
		v1.DELETE("/users/:id", h.User.Delete)

		v1.POST("/accounts", h.Account.Create)
		v1.GET("/accounts", h.Account.List)
		v1.GET("/accounts/:id", h.Account.GetByID)
		v1.PUT("/accounts/:id", h.Account.Update)
		v1.DELETE("/accounts/:id", h.Account.Delete)

		v1.POST("/transactions", h.Transaction.Create)
		v1.GET("/transactions", h.Transaction.List)
		v1.GET("/transactions/:id", h.Transaction.GetByID)
	}

	return r
}
