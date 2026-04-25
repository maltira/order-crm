package main

import (
	"log"
	"order-crm/config"
	"order-crm/internal/router"
	"order-crm/pkg/database"
)

func main() {
	config.InitEnv()
	db := database.InitDB()
	database.InitMongo()
	database.InitRedis()
	database.InitClickHouse()

	r := router.InitGinRouter(db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	log.Println("Server started at port 8080")
}
