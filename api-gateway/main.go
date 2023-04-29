package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"api-gateway/handlers"
	"api-gateway/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1 := r.Group("/v1")
	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{"ok": true})
	})

	users := v1.Group("/users")
	users.POST("/register", handlers.Register())
	users.POST("/login", handlers.Login())
	users.POST("/verify", handlers.Verify())

	products := v1.Group("/products")
	products.GET("/", handlers.GetAllProduct())
	products.GET("/:id", handlers.GetProduct())

	products.Use(middleware.VerifyToken)
	products.POST("/", handlers.CreateProduct())
	products.PATCH("/:id", handlers.UpdateProduct())
	products.DELETE("/:id", handlers.DeleteProduct())

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
