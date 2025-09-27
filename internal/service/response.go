package service

import (
	"fmt"
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"github.com/gin-gonic/gin"
)

func sendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"message": msg,
		"errCode": code,
	})
}

func sendSuccess(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("operation from handler 1231313: %s successfull", op),
		"data":    data,
	})

}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type CreateProductResponse struct {
	Message string                  `json:"message"`
	Data    schemas.ProductResponse `json:"data"`
}

type DeleteProductResponse struct {
	Message string                  `json:"message"`
	Data    schemas.ProductResponse `json:"data"`
}
type FindProductResponse struct {
	Message string                  `json:"message"`
	Data    schemas.ProductResponse `json:"data"`
}
type FindAllProductsResponse struct {
	Message string                    `json:"message"`
	Data    []schemas.ProductResponse `json:"data"`
}
type UpdateProductResponse struct {
	Message string                  `json:"message"`
	Data    schemas.ProductResponse `json:"data"`
}
