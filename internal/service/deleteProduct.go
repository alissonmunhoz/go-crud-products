package service

import (
	"fmt"
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"github.com/gin-gonic/gin"
)

// @BasePath /v1

// @Summary Delete product
// @Description Delete a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param id query string true "Product identification"
// @Success 200 {object} DeleteProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 504 {object} ErrorResponse
// @Router /product [delete]
func DeleteProductService(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}
	product := schemas.Product{}

	if err := db.First(&product, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, fmt.Sprintf("product with id: %s not found", id))
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error deleting product with id: %s", id))
		return
	}
	sendSuccess(ctx, "delete-product", product)
}
