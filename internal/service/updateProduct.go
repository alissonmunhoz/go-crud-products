// handler/update_product.go
package service

import (
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"

	"github.com/gin-gonic/gin"
)

// @BasePath /v1
// @Summary Update product
// @Description Update a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id query string true "Product identification"
// @Param request body UpdateProductRequest true "Product data to update"
// @Success 200 {object} UpdateProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /product [put]
func UpdateProductService(ctx *gin.Context) {
	var req UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("bind error: %v", err)
		sendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Errorf("validation error: %v", err)
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	var product schemas.Product
	if err := db.First(&product, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, "product not found")
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Quantity >= 0 {
		product.Quantity = req.Quantity
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if err := db.Save(&product).Error; err != nil {
		logger.Errorf("error updating product: %v", err)
		sendError(ctx, http.StatusInternalServerError, "error updating product")
		return
	}

	ctx.JSON(http.StatusOK, UpdateProductResponse{
		Message: "operation from handler: update-product successful",
		Data:    toProductResponse(product),
	})
}
