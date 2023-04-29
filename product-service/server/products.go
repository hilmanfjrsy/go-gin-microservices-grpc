package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"product-service/config"
	"product-service/models"
	"product-service/pb"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctx = context.Background()

type ProductServer struct {
	pb.UnimplementedProductServiceServer
}

func (a *ProductServer) getLatestVersion(version *uint64) error {
	err := config.DB.Unscoped().Model(&models.Products{}).Select("coalesce(max(version), 0) as version").Scan(version).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *ProductServer) getVersionById(version *uint64, id uint64) error {
	err := config.DB.Unscoped().Model(&models.Products{ID: uint(id)}).Select("coalesce(max(version), 0) as version").Scan(version).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *ProductServer) GetAll(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	productsModel := []models.Products{}
	products := &pb.GetProductsResponse{}

	var latestVersion uint64
	keyCacheVersion := "product_version"
	keyCacheProduct := "products"
	err := a.getLatestVersion(&latestVersion)
	if err != nil {
		log.Println("Failed to get latest version:", err.Error())
		return nil, err
	}

	cacheVersion := 0
	version, err := config.Redis.Get(ctx, keyCacheVersion).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Failed to get cache version  fromredis:", err.Error())
			return nil, err
		}
	} else {
		cacheVersion, err = strconv.Atoi(version)
		if err != nil {
			log.Println("Failed to convert string to int:", err.Error())
			return nil, err
		}
	}

	cacheProduct, err := config.Redis.Get(ctx, keyCacheProduct).Result()
	if err == redis.Nil || uint64(cacheVersion) < latestVersion {
		log.Println("Get product from database")
		pagination := models.Pagination{CurrentPage: req.Page}
		var count int64
		offset := (req.Page - 1) * req.Limit
		err := config.DB.Model(&models.Products{}).Count(&count).Error
		if err != nil {
			log.Println("Failed to get total product:", err.Error())
			return nil, err
		}
		pagination.TotalPage = uint64(math.Ceil(float64(count) / float64(req.Limit)))
		if req.Page < uint64(pagination.TotalPage) {
			pagination.NextPage = req.Page + 1
		}
		if req.Page > 1 {
			pagination.PrevPage = req.Page - 1
		}

		err = config.DB.Model(&models.Products{}).Offset(int(offset)).Limit(int(req.Limit)).Find(&productsModel).Error
		if err != nil {
			log.Println("Failed to get product:", err.Error())
			return nil, err
		}
		b, _ := json.Marshal(productsModel)
		p, _ := json.Marshal(pagination)
		cacheProduct = fmt.Sprintf(`{"result":%s,"pagination":%s}`, b, p)

		err = config.Redis.Set(ctx, keyCacheVersion, latestVersion, 0).Err()
		if err != nil {
			log.Println("Failed to set cache version to redis:", err.Error())
			return nil, err
		}

		err = config.Redis.Set(ctx, keyCacheProduct, cacheProduct, 0).Err()
		if err != nil {
			log.Println("Failed to set cache product to redis:", err.Error())
			return nil, err
		}
	} else if err != nil {
		log.Println("Failed to get cache from redis:", err.Error())
		return nil, err
	} else {
		log.Println("Get product from cache")
	}

	err = protojson.Unmarshal([]byte(cacheProduct), products)
	if err != nil {
		log.Println("Failed to unmarshal product:", err.Error())
		return nil, err
	}
	return products, nil
}

func (a *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	productsModel := models.Products{}
	product := &pb.GetProductResponse{}

	var latestVersion uint64
	keyCacheVersionProduct := fmt.Sprintf("product_%v_version", req.Id)
	keyCacheProduct := fmt.Sprintf("product_%v", req.Id)
	err := a.getVersionById(&latestVersion, req.Id)
	if err != nil {
		log.Println("Failed to get version by id:", err.Error())
		return nil, err
	}

	cacheVersion := 0
	version, err := config.Redis.Get(ctx, keyCacheVersionProduct).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Failed to get cache version from redis:", err.Error())
			return nil, err
		}
	} else {
		cacheVersion, err = strconv.Atoi(version)
		if err != nil {
			log.Println("Failed to convert string to int:", err.Error())
			return nil, err
		}
	}

	cacheProduct, err := config.Redis.Get(ctx, keyCacheProduct).Result()
	if err == redis.Nil || uint64(cacheVersion) < latestVersion {
		err := config.DB.First(&productsModel, req.Id).Error
		if err != nil {
			log.Println("Failed to get product:", err.Error())
			return nil, err
		}
		b, _ := json.Marshal(productsModel)
		cacheProduct = fmt.Sprintf(`{"result":%s}`, b)

		err = config.Redis.Set(ctx, keyCacheVersionProduct, latestVersion, 0).Err()
		if err != nil {
			log.Println("Failed to set cache version to redis:", err.Error())
			return nil, err
		}

		err = config.Redis.Set(ctx, keyCacheProduct, cacheProduct, 0).Err()
		if err != nil {
			log.Println("Failed to set cache product to redis:", err.Error())
			return nil, err
		}
	} else if err != nil {
		log.Println("Failed to get cache from redis:", err.Error())
		return nil, err
	} else {
		log.Println("Get product from cache")
	}
	log.Println("Latest cache version:", latestVersion, "-> Cache version:", cacheVersion)

	err = jsonpb.UnmarshalString(cacheProduct, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if req.Name == "" || req.Price == 0 || req.UserId == 0 {
		return nil, errors.New("Missing required fields")
	}

	var latestVersion uint64
	err := a.getLatestVersion(&latestVersion)
	if err != nil {
		log.Println("Failed to get latest version:", err.Error())
		return nil, err
	}

	product := &models.Products{
		Name:        req.Name,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		UserId:      req.UserId,
		Price:       req.Price,
		Version:     latestVersion + 1,
	}
	err = config.DB.Model(&models.Products{}).Create(&product).Error
	if err != nil {
		log.Println("Failed to create product:", err.Error())
		return nil, err
	}

	return &pb.CreateProductResponse{Result: "Product created!"}, nil
}

func (a *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	if req.Name == "" || req.Price == 0 || req.UserId == 0 {
		return nil, errors.New("Missing required fields")
	}

	var latestVersion uint64
	err := a.getLatestVersion(&latestVersion)
	if err != nil {
		log.Println("Failed to get latest version:", err.Error())
		return nil, err
	}

	product := &models.Products{
		Name:        req.Name,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Price:       req.Price,
		Version:     latestVersion + 1,
	}
	err = config.DB.Where("id = ? and user_id = ?", req.Id, req.UserId).Updates(product).Error
	if err != nil {
		log.Println("Failed to update product:", err.Error())
		return nil, err
	}

	return &pb.UpdateProductResponse{Result: "Product updated!"}, nil
}

func (a *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	product := models.Products{}
	err := config.DB.Where("user_id = ?", req.UserId).First(&product, req.Id).Error
	if err != nil {
		log.Println("Failed to delete product:", err.Error())
		return nil, err
	}

	var latestVersion uint64
	err = a.getLatestVersion(&latestVersion)
	if err != nil {
		log.Println("Failed to get latest version:", err.Error())
		return nil, err
	}
	product.Version = latestVersion + 1
	err = config.DB.Where("id = ? and user_id = ?", req.Id, req.UserId).Updates(product).Error
	if err != nil {
		log.Println("Failed to update product:", err.Error())
		return nil, err
	}

	err = config.DB.Where("user_id = ?", req.UserId).Delete(&product, req.Id).Error
	if err != nil {
		log.Println("Failed to delete product:", err.Error())
		return nil, err
	}

	return &pb.DeleteProductResponse{Result: "Product deleted!"}, nil
}
