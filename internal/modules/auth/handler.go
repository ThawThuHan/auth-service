package auth

import (
	"net/http"

	"github.com/TRIOSYS-Software/auth-service/internal/config"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/dto"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/service"
	"github.com/TRIOSYS-Software/auth-service/internal/util"
	"github.com/TRIOSYS-Software/auth-service/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	log := config.Cfg.Logger
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid request payload", zap.Error(err))
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Error("Authentication failed", zap.Error(err))
		util.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed")
		return
	}
	util.SuccessResponse(c, http.StatusOK, token)
}

func (h *AuthHandler) JWKS(c *gin.Context) {
	log := config.Cfg.Logger
	jwks, err := jwt.ServeJWKS()
	if err != nil {
		log.Error("Failed to serve JWKS", zap.Error(err))
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to serve JWKS")
		return
	}
	util.SuccessResponse(c, http.StatusOK, jwks)
}
