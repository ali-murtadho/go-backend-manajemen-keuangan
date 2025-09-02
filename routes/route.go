package routes

import (
	"backend-api/controllers"
	"backend-api/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// route register
	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	router.GET("/api/users", middlewares.AuthMiddleware(), controllers.FindUsers)
	router.POST("/api/users", middlewares.AuthMiddleware(), controllers.CreateUser)
	router.GET("/api/users/:id", middlewares.AuthMiddleware(), controllers.FindUserById)
	router.PUT("/api/users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
	router.DELETE("/api/users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)

	router.GET("/api/times", middlewares.AuthMiddleware(), controllers.GetListWaktu)
	router.GET("/api/times/:bulan_tahun", middlewares.AuthMiddleware(), controllers.GetWaktuByBulanTahun)
	router.POST("/api/times", middlewares.AuthMiddleware(), controllers.CreateTimes)
	router.PUT("/api/times/:id", middlewares.AuthMiddleware(), controllers.UpdateTimes)
	router.DELETE("/api/times/:id", middlewares.AuthMiddleware(), controllers.DeleteTimes)
	return router
}
