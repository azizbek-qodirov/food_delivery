package api

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gateway-admin/api/docs"
	"gateway-admin/api/handlers"
	"gateway-admin/api/middleware"
)

// @title Swaggers of admin panel
// @version 1.0
// @BasePath /api/swagger/index.html#/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(h *handlers.HTTPHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := router.Group("/", middleware.JWTMiddleware())
	protected.Use(middleware.IsAdminMiddleware())

	protected.POST("/add-product", h.AddProduct)
	protected.DELETE("/delete-product/:id", h.DeleteProduct)

	return router
}
