package main

import (
	"backend-api/config"
	"backend-api/database"
	"backend-api/helpers"
	"backend-api/routes"
)

func main() {

	helpers.Init() // register the "pwd" rule

	//load config .env
	config.LoadEnv()

	database.InitDB()

	//setup router
	r := routes.SetupRouter()

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
