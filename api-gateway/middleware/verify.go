package middleware

import (
	"api-gateway/pb"
	"api-gateway/utils"
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyToken(ctx *gin.Context) {
	token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		log.Println("Token required")
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	conn, err := utils.GRPCClient(os.Getenv("GRPC_AUTH_HOST"))
	if err != nil {
		log.Println("Error Connection to GRPC", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Verify(c, &pb.UserVerifyRequest{
		Token: token,
	})
	if err != nil {
		log.Println("Error Verify", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	if !response.Valid {
		log.Println("Error ", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	ctx.Set("user_id", response.UserId)
	ctx.Next()
}
