package main

import (
	"backend-api/config"
	"backend-api/database"

	"github.com/gin-gonic/gin"
)

func main() {

	//load config .env
	config.LoadEnv()

	database.InitDB()

	//inisialiasai Gin
	router := gin.Default()

	//membuat route dengan method GET
	router.GET("/", func(c *gin.Context) {

		//return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World! Hello Gin go!",
		})
	})

	//mulai server
	router.Run(":" + config.GetEnv("PORT", "3000"))
}
