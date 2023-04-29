package main

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/pb"
	"auth-service/server"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	loadEnv()
	r := gin.Default()
	log.Println("Starting auth service")
	err := config.ConnectToPostgres()
	if err != nil {
		log.Panic("Can't connect to Postgres!", err)
	}
	models.AutoMigrateModels()

	go listenGRPC()

	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Panic(err)
	}
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

	pb.RegisterUserServiceServer(s, &server.AuthServer{})

	log.Printf("gRPC Auth Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC Auth: %v", err)
	}
}
