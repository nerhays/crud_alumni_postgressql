package main

import (
	"crud_alumni/config"
	"crud_alumni/database"
	"log"
)

func main() {
	config.LoadEnv()
	config.InitLogger()
	database.ConnectDB()

	app := config.App()

	// ✅ port ambil dari .env
	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}

