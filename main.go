package main

import (
	"backend-api/config"
	"backend-api/database"
	"backend-api/helpers"
	"backend-api/routes"
	"fmt"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend-api/docs" // ini wajib, nanti setelah "swag init"
)

// @title Backend API
// @version 1.0
// @description API documentation untuk backend golang
// @host localhost:3000
// @BasePath /api
func main() {
	helpers.Init() // register the "pwd" rule

	//load config .env
	config.LoadEnv()

	database.InitDB()

	//setup router
	r := routes.SetupRouter()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ambil port dari env (default 3000)
	port := config.GetEnv("APP_PORT", "3000")

	// tampilkan info lebih lengkap di console
	fmt.Printf("🚀 Server running at: http://localhost:%s/swagger/index.html\n", port)

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
