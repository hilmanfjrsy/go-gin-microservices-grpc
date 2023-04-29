package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"product-service/config"
	"product-service/models"
	"product-service/pb"
	"product-service/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	loadEnv()
	r := gin.Default()
	err := config.ConnectToMySQL()
	if err != nil {
		log.Panic("Can't connect to MySQL!", err)
	}
	models.AutoMigrateModels()

	err = config.ConnectToRedis()
	if err != nil {
		log.Panic("Can't connect to Redis!", err)
	}

	go listenGRPC()

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func listenGRPC() {
	gRpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC Auth: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, &server.ProductServer{})

	log.Printf("gRPC Product Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC Auth: %v", err)
	}
}
