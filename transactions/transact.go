package transactions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	Transact(*gin.Context)
}

type impl struct {
	UseCase
}

func NewHandler(useCase UseCase) Handler {
	return &impl{
		UseCase: useCase,
	}
}

func (i *impl) Transact(ctx *gin.Context) {
	var transaction Transaction

	fmt.Println("here.")
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println("here..")
	if err := i.UseCase.Transact(transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println("here...")

	ctx.JSON(http.StatusOK, "transaction processed successfully")
}
