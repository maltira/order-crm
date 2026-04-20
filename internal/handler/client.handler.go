package handler

import (
	"net/http"
	"order-crm/internal/model/dto"
	"order-crm/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	sc service.ClientService
}

func NewClientHandler(sc service.ClientService) *ClientHandler {
	return &ClientHandler{sc: sc}
}

func (h *ClientHandler) GetAllClients(c *gin.Context) {
	clients, err := h.sc.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

func (h *ClientHandler) GetClientById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	client, err := h.sc.GetClientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req dto.ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.sc.CreateClient(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Клиент успешно создан",
		"client":  client,
	})
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	var req dto.ClientRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.sc.UpdateClient(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Клиент успешно обновлён"})
}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}
	err = h.sc.DeleteClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Клиент удалён"})
}
