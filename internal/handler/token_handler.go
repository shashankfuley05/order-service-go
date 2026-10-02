package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/dto"
	"github.com/shashank/order-service/internal/service"
)

type TokenHandler struct {
	authService service.AuthService
}

func NewTokenHandler(service service.AuthService) *TokenHandler {

	return &TokenHandler{
		authService: service,
	}
}

func (handler *TokenHandler) GenerateTokenForUser(ctx *gin.Context) {

	var tokenRequest dto.TokenRequest
	err := ctx.ShouldBindBodyWithJSON(&tokenRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	token, err := handler.authService.GenerateToken(&tokenRequest)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, token)

}
