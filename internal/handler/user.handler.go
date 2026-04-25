package handler

import (
	"net/http"
	"order-crm/internal/model/dto"
	"order-crm/internal/service"
	"order-crm/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	sc service.UserService
}

func NewUserHandler(sc service.UserService) *UserHandler {
	return &UserHandler{
		sc: sc,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req *dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.sc.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go database.LogUserEvent(userID, "create_user")

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	user, err := h.sc.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.sc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.sc.UpdateUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go database.LogUserEvent(userID, "update_user")

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь обновлён"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	err = h.sc.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	go database.LogUserEvent(userID, "delete_user")

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удалён"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.sc.ChangePassword(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go database.LogUserEvent(userID, "change_user_password")

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменён"})
}
