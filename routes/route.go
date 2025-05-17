package routes

import (
	"backend-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	// route register
	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	return router
}
