// handler/find_all_products.go
package service

import (
	"net/http"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"github.com/gin-gonic/gin"
)

// @BasePath /v1
// @Summary Find All products
// @Description Find all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} FindAllProductsResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func FindAllProductsService(ctx *gin.Context) {
	var products []schemas.Product
	if err := db.Find(&products).Error; err != nil {
		sendError(ctx, http.StatusInternalServerError, "error listing products")
		return
	}

	resp := make([]schemas.ProductResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, toProductResponse(p))
	}

	ctx.JSON(http.StatusOK, FindAllProductsResponse{
		Message: "operation from handler: list-products successful",
		Data:    resp,
	})
}
