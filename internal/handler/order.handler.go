package handler

import (
	"net/http"
	"order-crm/internal/model/dto"
	"order-crm/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	sc service.OrderService
}

func NewOrderHandler(sc service.OrderService) *OrderHandler {
	return &OrderHandler{sc: sc}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.sc.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	order, err := h.sc.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.sc.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.sc.UpdateOrderStatus(id, req.IDStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус заказа успешно обновлён"})
}

func (h *OrderHandler) AddPayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	var req dto.AddPaymentRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.sc.AddPaymentToOrder(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Платёж успешно добавлен"})
}

func (h *OrderHandler) AddOrderItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	var req dto.OrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.sc.AddOrderItem(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Позиция успешно добавлена в заказ"})
}

func (h *OrderHandler) DeleteOrderItem(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil || orderId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}
	orderItemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || orderItemId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	if err = h.sc.DeleteOrderItem(orderId, orderItemId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Позиция успешно удалена"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	if err = h.sc.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заказ успешно удалён"})
}
