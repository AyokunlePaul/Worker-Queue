package main

import (
	"LemFi/authentication"
	"LemFi/transactions"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	authUseCase := authentication.NewUseCase()
	transactionUseCase := transactions.NewUseCase(authUseCase)

	authHandler := authentication.NewAuthHandler(authUseCase)
	transactionHandler := transactions.NewHandler(transactionUseCase)

	totalVerificationWorkers := 2
	totalTransactionWorkers := 2
	verificationSleepDuration := 30 * time.Second
	transactionSleepDuration := 20 * time.Second

	// Create workers
	for i := 0; i < totalVerificationWorkers; i++ {
		i := i
		go func() {
			for {
				fmt.Printf("running verification worker %d\n", i)
				_ = authUseCase.VerifyUser(i)
				time.Sleep(verificationSleepDuration)
			}
		}()
	}

	// Create transaction workers
	for i := 0; i < totalTransactionWorkers; i++ {
		i := i
		go func() {
			for {
				fmt.Printf("running transaction worker %d\n", i)
				_ = transactionUseCase.ProcessTransaction(i)
				time.Sleep(transactionSleepDuration)
			}
		}()
	}

	engine := gin.Default()
	engine.POST("/create", authHandler.CreateAccount)
	engine.GET("/users", authHandler.GetAllUsers)
	engine.POST("/transact", transactionHandler.Transact)

	err := engine.Run()
	if err != nil {
		fmt.Println(err)
	}
}
