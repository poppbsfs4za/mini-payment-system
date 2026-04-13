package main

import (
	"log"

	"mini-payment-system/internal/config"
	"mini-payment-system/internal/database"
	"mini-payment-system/internal/delivery/http/handlers"
	"mini-payment-system/internal/delivery/http/router"
	"mini-payment-system/internal/domain/entities"
	gormrepo "mini-payment-system/internal/infrastructure/persistence/gorm"
	accountuc "mini-payment-system/internal/usecase/account"
	transactionuc "mini-payment-system/internal/usecase/transaction"
	useruc "mini-payment-system/internal/usecase/user"

	_ "mini-payment-system/docs/swagger"
)

// @title Mini Payment System API
// @version 1.0
// @description Clean Architecture CRUD API assignment for Arise by Infinitas.
// @BasePath /api/v1
func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&entities.User{}, &entities.Account{}, &entities.Transaction{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userRepo := gormrepo.NewUserRepository(db)
	accountRepo := gormrepo.NewAccountRepository(db)
	transactionRepo := gormrepo.NewTransactionRepository(db)
	txManager := gormrepo.NewTxManager(db)

	userService := useruc.NewService(userRepo)
	accountService := accountuc.NewService(accountRepo, userRepo)
	transactionService := transactionuc.NewService(accountRepo, transactionRepo, txManager)

	r := router.New(router.Handlers{
		User:        handlers.NewUserHandler(userService),
		Account:     handlers.NewAccountHandler(accountService),
		Transaction: handlers.NewTransactionHandler(transactionService),
	})

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
