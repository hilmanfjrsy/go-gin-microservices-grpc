package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"api-gateway/pb"
	"api-gateway/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type CreateProductData struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Thumbnail   string  `json:"thumbnail"`
}

type UpdateProductData struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Thumbnail   string  `json:"thumbnail"`
}

func CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createProduct := CreateProductData{}
		userId := ctx.GetUint64("user_id")

		if err := ctx.ShouldBindJSON(&createProduct); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_PRODUCT_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewProductServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.CreateProduct(c, &pb.CreateProductRequest{
			Name:        createProduct.Name,
			Description: createProduct.Description,
			Price:       createProduct.Price,
			Thumbnail:   createProduct.Thumbnail,
			UserId:      userId,
		})
		if err != nil {
			log.Println("Failed to create", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		updateProduct := UpdateProductData{}
		userId := ctx.GetUint64("user_id")

		if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_PRODUCT_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewProductServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.UpdateProduct(c, &pb.UpdateProductRequest{
			Id:          uint64(id),
			Name:        updateProduct.Name,
			Description: updateProduct.Description,
			Price:       updateProduct.Price,
			Thumbnail:   updateProduct.Thumbnail,
			UserId:      userId,
		})
		if err != nil {
			log.Println("Failed to update", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		conn, err := utils.GRPCClient(os.Getenv("GRPC_PRODUCT_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewProductServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.DeleteProduct(c, &pb.DeleteProductRequest{
			Id:     uint64(id),
			UserId: ctx.GetUint64("user_id"),
		})
		if err != nil {
			log.Println("Failed to delete", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func GetAllProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p := ctx.Query("page")
		l := ctx.Query("limit")
		if p == "" {
			p = "1"
		}
		if l == "" {
			l = "10"
		}
		page, err := strconv.Atoi(p)
		if err != nil {
			log.Println("Failed to convert query page", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		limit, err := strconv.Atoi(l)
		if err != nil {
			log.Println("Failed to convert query limit", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		conn, err := utils.GRPCClient(os.Getenv("GRPC_PRODUCT_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewProductServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.GetAll(c, &pb.GetProductsRequest{
			Page:  uint64(page),
			Limit: uint64(limit),
		})
		if err != nil {
			log.Println("Failed to get all", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		b, err := protojson.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		products := map[string]any{}
		err = json.Unmarshal(b, &products)

		if products["result"] == nil {
			products["result"] = []interface{}{}
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, products)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		conn, err := utils.GRPCClient(os.Getenv("GRPC_PRODUCT_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewProductServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.GetProduct(c, &pb.GetProductRequest{Id: uint64(id)})
		if err != nil {
			log.Println("Failed to get product", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		b, err := protojson.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		products := map[string]interface{}{}
		err = json.Unmarshal(b, &products)

		utils.ResponseSuccess(ctx, http.StatusAccepted, products)
	}
}
