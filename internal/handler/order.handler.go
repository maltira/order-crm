package handler

import (
	"log"
	"net/http"
	"order-crm/internal/model/dto"
	"order-crm/internal/service"
	"order-crm/pkg/database"
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
	userID := c.MustGet("user_id").(int)
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

	go database.LogUserEvent(userID, "create_order")

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
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

	go database.LogUserEvent(userID, "update_order_status")
	go database.LogOrderStatus(id, req.IDStatus)

	order, err := h.sc.GetOrderById(id)
	if err != nil {
		log.Println("Ошибка получения заказа для redis", err.Error())
	} else {
		go database.IncrPurchases(order.Items)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус заказа успешно обновлён"})
}

func (h *OrderHandler) AddPayment(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
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

	go database.LogUserEvent(userID, "add_payment")

	c.JSON(http.StatusOK, gin.H{"message": "Платёж успешно добавлен"})
}

func (h *OrderHandler) AddOrderItem(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
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

	go database.LogUserEvent(userID, "add_order_item")

	c.JSON(http.StatusOK, gin.H{"message": "Позиция успешно добавлена в заказ"})
}

func (h *OrderHandler) DeleteOrderItem(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
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

	go database.LogUserEvent(userID, "delete_order_item")

	c.JSON(http.StatusOK, gin.H{"message": "Позиция успешно удалена"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Передан неверный формат ID"})
		return
	}

	if err = h.sc.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go database.LogUserEvent(userID, "delete_order")

	c.JSON(http.StatusOK, gin.H{"message": "Заказ успешно удалён"})
}
