package database

import (
	"context"
	"fmt"
	"log"
	"order-crm/internal/model"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection error:", err)
	}

	log.Println("Redis connected")
}

func IncrPurchases(items []model.OrderItem) {
	for _, el := range items {
		key := fmt.Sprintf("product:%d:purchases", el.ProductID)
		err := RDB.IncrBy(context.Background(), key, 1).Err()
		if err != nil {
			log.Println("Redis incr error:", err)
		}
	}
}
