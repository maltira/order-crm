package router

import (
	"order-crm/internal/handler"
	"order-crm/internal/repository"
	"order-crm/internal/service"
	"order-crm/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitGinRouter() *gin.Engine {
	uR := repository.NewUserRepository(database.GetDB())
	uSc := service.NewUserService(uR)
	uHd := handler.NewUserHandler(uSc)

	r := gin.Default()
	api := r.Group("/api")

	user := api.Group("/user")

	{
		user.POST("", uHd.CreateUser)

		user.GET("/id/:id", uHd.GetUserById)
		user.GET("/login/:login", uHd.GetUserByLogin)
		user.GET("/all", uHd.GetAllUsers)

		user.PUT("", uHd.UpdateUser)
		user.DELETE("/:id", uHd.DeleteUser)
		user.PUT("/password", uHd.ChangePassword)
	}

	return r
}
