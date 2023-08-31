package authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	CreateAccount(*gin.Context)
	GetAllUsers(*gin.Context)
}

type auth struct {
	UseCase
}

func NewAuthHandler(authUseCase UseCase) Handler {
	return &auth{
		UseCase: authUseCase,
	}
}

func (a *auth) CreateAccount(ctx *gin.Context) {
	var err error
	name := ctx.PostForm("name")

	user := &User{
		Name:     name,
		Balance:  1000,
		Verified: false,
	}

	if err = a.CreateUser(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, "user created successfully")
}

func (a *auth) GetAllUsers(ctx *gin.Context) {
	users, _ := a.UseCase.GetAllUsers()

	ctx.JSON(http.StatusOK, users)
}
