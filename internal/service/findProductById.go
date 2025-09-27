package service

import (
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"github.com/gin-gonic/gin"
)

// @BasePath /v1

// @Summary Find product
// @Description Find a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id query string true "Product identification"
// @Success 200 {object} FindProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /product [get]
func FindProductService(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}
	product := schemas.Product{}
	if err := db.First(&product, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, "product not found")
		return
	}

	sendSuccess(ctx, "show-product", product)
}
