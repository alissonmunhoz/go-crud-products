package router

import (
	_ "github.com/alissonmunhoz/go-crud-products/docs"
	service "github.com/alissonmunhoz/go-crud-products/internal/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeRoutes(router *gin.Engine) {
	service.InitializeHandler()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := router.Group("/v1")

	{
		v1.POST("/product", service.CreateProductService)
		v1.DELETE("/product", service.DeleteProductService)
		v1.PUT("/product", service.UpdateProductService)
		v1.GET("/products", service.FindAllProductsService)
		v1.GET("/product", service.FindProductService)
	}

}
