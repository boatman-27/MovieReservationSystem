package main

import (
	"log"
	DB "movie/config"
	"movie/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	DB := DB.ConnectDB()
	routes.SetupRouter(router, DB)

	router.Run(":8000")
}
