package routes

import (
	controller "github.com/dev-hack95/jwt-golang/controllers"
	middleware "github.com/dev-hack95/jwt-golang/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// User will have token when he loged in or sign in so we make sure that the user will pe protected using middleware
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())
}
