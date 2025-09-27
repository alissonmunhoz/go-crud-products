// handler/create_product.go
package service

import (
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"

	"github.com/gin-gonic/gin"
)

// @BasePath /v1
// @Summary Create product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Request body"
// @Success 200 {object} CreateProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /product [post]
func CreateProductService(ctx *gin.Context) {
	var req CreateProductRequest
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

	product := schemas.Product{
		Name:        req.Name,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Description: req.Description,
	}

	if err := db.Create(&product).Error; err != nil {
		logger.Errorf("error creating product: %v", err)
		sendError(ctx, http.StatusInternalServerError, "error creating product on database")
		return
	}

	ctx.JSON(http.StatusOK, CreateProductResponse{
		Message: "operation from handler: create-product successful",
		Data:    toProductResponse(product),
	})
}
