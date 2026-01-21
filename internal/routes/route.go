package routes

import (
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *auth.AuthHandler) {
	authGroup := router.Group("/auth")
	{
		setupAuthRoutes(authGroup, authHandler)
	}
}

func setupAuthRoutes(route *gin.RouterGroup, authHandler *auth.AuthHandler) {
	route.POST("/login", authHandler.Login)
	route.GET("/.well-known/jwks.json", authHandler.JWKS)
}
