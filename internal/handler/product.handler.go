package handler

import (
	"context"
	"errors"
	"fmt"
	mng "order-crm/internal/repository/mongo"
	"order-crm/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetAllProductsHandler(c *gin.Context) {
	products, err := mng.GetAllProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}

func GetProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Неверный формат ID"})
		return
	}
	product, err := mng.GetProductById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, product)
}

func GetProductPurchases(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Неверный формат ID"})
		return
	}

	key := fmt.Sprintf("product:%d:purchases", id)
	val, err := database.RDB.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(200, 0)
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	purchases, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка преобразования данных " + err.Error()})
		return
	}

	c.JSON(200, purchases)
}
