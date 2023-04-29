package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"api-gateway/pb"
	"api-gateway/utils"

	"github.com/gin-gonic/gin"
)

type RegisterUser struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUser struct {
	Token string `json:"token" binding:"required"`
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		registerUser := RegisterUser{}

		if err := ctx.ShouldBindJSON(&registerUser); err != nil {
			log.Println("Failed to binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_AUTH_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.Register(c, &pb.UserRegisterRequest{
			Name:     registerUser.Name,
			Username: registerUser.Username,
			Password: registerUser.Password,
		})

		if err != nil {
			log.Println("Failed to register", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loginUser := LoginUser{}

		if err := ctx.ShouldBindJSON(&loginUser); err != nil {
			log.Println("Failed to binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_AUTH_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.Login(c, &pb.UserLoginRequest{
			Username: loginUser.Username,
			Password: loginUser.Password,
		})

		if err != nil {
			log.Println("Failed to login", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyUser := VerifyUser{}
		if err := ctx.ShouldBindJSON(&verifyUser); err != nil {
			log.Println("Failed to binding", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_AUTH_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.Verify(c, &pb.UserVerifyRequest{
			Token: verifyUser.Token,
		})
		if err != nil {
			log.Println("Failed to verify", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}
