package router

import (
	"database/sql"
	"order-crm/internal/handler"
	"order-crm/internal/middleware"
	"order-crm/internal/repository"
	"order-crm/internal/service"

	"github.com/gin-gonic/gin"
)

func InitGinRouter(db *sql.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewClientRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userService := service.NewUserService(userRepo)
	clientService := service.NewClientService(clientRepo)
	orderService := service.NewOrderService(orderRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)
	orderHandler := handler.NewOrderHandler(orderService)

	r := gin.Default()
	r.Use(middleware.CORS())
	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		auth.POST("/refresh", authHandler.Refresh)
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
	}

	users := api.Group("/users").Use(middleware.AuthMiddleware()).Use(middleware.AdminOnly())
	{
		users.POST("", userHandler.CreateUser)
		users.PUT("", userHandler.UpdateUser)
		users.GET("", userHandler.GetAllUsers)

		users.GET("/:id", userHandler.GetUserById)
		users.DELETE("/:id", userHandler.DeleteUser)

		users.PUT("/change-password", userHandler.ChangePassword)
	}

	clients := api.Group("/clients").Use(middleware.AuthMiddleware())
	{
		clients.GET("", clientHandler.GetAllClients)
		clients.POST("", middleware.ManagerOrHigher(), clientHandler.CreateClient)

		clients.PUT("/:id", middleware.ManagerOrHigher(), clientHandler.UpdateClient)
		clients.GET("/:id", clientHandler.GetClientById)
		clients.DELETE("/:id", middleware.ManagerOrHigher(), clientHandler.DeleteClient)
	}

	orders := api.Group("/orders").Use(middleware.AuthMiddleware())
	{
		orders.GET("", orderHandler.GetAllOrders)
		orders.GET("/:id", orderHandler.GetOrderById)
		orders.POST("", orderHandler.CreateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)

		orders.POST("/:id/items", orderHandler.AddOrderItem)
		orders.DELETE("/:id/items/:item_id", orderHandler.DeleteOrderItem)

		orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orders.POST("/:id/payments", middleware.ManagerOrHigher(), orderHandler.AddPayment)

		// products
		orders.GET("/products", handler.GetAllProductsHandler)
		orders.GET("/products/:id", handler.GetProductById)
		orders.GET("/products/:id/purchases", handler.GetProductPurchases)
	}

	return r
}
