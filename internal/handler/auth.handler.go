package handler

import (
	"log"
	"net/http"
	"order-crm/internal/model/dto"
	"order-crm/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	sc service.UserService
}

func NewAuthHandler(sc service.UserService) *AuthHandler {
	return &AuthHandler{sc: sc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.sc.Login(req.Login, req.Password)
	if err != nil {
		switch err.Error() {
		case "user is blocked":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь заблокирован администратором"})
		case "invalid credentials":
			c.JSON(http.StatusNotFound, gin.H{"error": "Неверный логин или пароль"})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	access, refresh, err := h.sc.GenerateAndSaveTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// рефреш токен в куки хттп онли
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: access, // access в state management
		User:        *user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	if refreshToken != "" {
		if err := h.sc.RevokeRefreshToken(refreshToken); err != nil {
			log.Println("Revoke refresh token error:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Вы вышли из аккаунта",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.sc.RefreshToken(refreshToken)
	if err != nil {
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)

		switch err.Error() {
		case "invalid or expired refresh token":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ваша сессия больше недействительна"})
		case "user is blocked":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ваш аккаунт заблокирован"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	access, refresh, err := h.sc.GenerateAndSaveTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: access, // access в state management
		User:        *user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	user, err := h.sc.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
