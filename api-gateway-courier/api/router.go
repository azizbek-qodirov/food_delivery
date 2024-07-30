package api

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/api/docs"
	"auth-service/api/handlers"
	"auth-service/api/middleware"
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

	protected.PUT("/ban/:id", h.BanUser)
	protected.PUT("/unban/:id", h.UnbanUser)

	protected.POST("/add-courier", h.AddCourier)
	protected.DELETE("/delete-courier/:id", h.DeleteCourier)

	protected.POST("/add-product-manager", h.AddProductManager)
	protected.DELETE("/delete-product-manager/:id", h.DeleteProductManager)

	return router
}
