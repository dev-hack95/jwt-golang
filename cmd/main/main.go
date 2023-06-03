package main

import (
	"os"

	routes "github.com/dev-hack95/jwt-golang/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthenticationRoutes(router)
	routes.UserRoutes(router)

	//func (ctx *gin.Context)  {} == func(w http.ResponseWriter, r *http.Request) {}

	router.Run(":" + port)
}
