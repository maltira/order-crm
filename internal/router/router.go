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
	uR := repository.NewUserRepository(database.GetDB())
	uSc := service.NewUserService(uR)
	uHd := handler.NewUserHandler(uSc)
	authHandler := handler.NewAuthHandler(uSc)

	r := gin.Default()
	api := r.Group("/api")

	user := api.Group("/user").Use(middleware.AuthMiddleware())
	{
		user.POST("", uHd.CreateUser)

		user.GET("/id/:id", uHd.GetUserById)
		user.GET("/login/:login", uHd.GetUserByLogin)
		user.GET("/all", uHd.GetAllUsers)

		user.PUT("", uHd.UpdateUser)
		user.DELETE("/:id", uHd.DeleteUser)
		user.PUT("/password", uHd.ChangePassword)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.Refresh)
	}

	return r
}
