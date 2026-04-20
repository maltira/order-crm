package router

import (
	"order-crm/internal/handler"
	"order-crm/internal/middleware"
	"order-crm/internal/repository"
	"order-crm/internal/service"
	"order-crm/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitGinRouter() *gin.Engine {
	userRepo := repository.NewUserRepository(database.GetDB())
	clientRepo := repository.NewClientRepository(database.GetDB())
	userService := service.NewUserService(userRepo)
	clientService := service.NewClientService(clientRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)

	r := gin.Default()
	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.Refresh)
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

	return r
}
