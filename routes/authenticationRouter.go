package routes

import (
	controller "github.com/dev-hack95/jwt-golang/controllers"
	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(incomingRoutes *gin.Engine) {
	// User doesnt have token at this point of time
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())
}
