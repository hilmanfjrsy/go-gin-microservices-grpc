package server

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/pb"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthServer struct {
	pb.UnimplementedUserServiceServer
}

type AuthClaims struct {
	UserId uint64 `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

func (a *AuthServer) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	if req.Name == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("Name, username and password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err.Error())
		return nil, err
	}

	userEntry := &models.Users{
		Name:      req.Name,
		Username:  strings.ToLower(req.Username),
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = config.DB.Model(&models.Users{}).Create(userEntry).Error
	if err != nil {
		log.Println("Failed to register user:", err.Error())
		return nil, err
	}

	return &pb.UserRegisterResponse{Result: "Users created!"}, nil
}

func (a *AuthServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("Username and password is required")
	}

	user := &models.Users{}
	err := config.DB.Model(&models.Users{}).Where("username = ?", strings.ToLower(req.Username)).First(user).Error
	if err != nil {
		log.Println("Failed to get users by username:", err.Error())
		return nil, errors.New("Invalid username/password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("Failed to compare password:", err.Error())
		return nil, errors.New("Invalid username/password")
	}

	token, err := generateToken(user)
	if err != nil {
		log.Println("Failed to generate token:", err.Error())
		return nil, errors.New("Invalid username/password")
	}

	return &pb.UserLoginResponse{Token: token}, nil
}

func (a *AuthServer) Verify(ctx context.Context, req *pb.UserVerifyRequest) (*pb.UserVerifyResponse, error) {
	if req.Token == "" {
		return &pb.UserVerifyResponse{Valid: false, Message: "Token is required"}, errors.New("Token is required")
	}

	claims, err := parseToken(req.Token, os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Println("Failed to parse token:", err.Error())
		return &pb.UserVerifyResponse{Valid: false, Message: err.Error()}, err
	}

	user := &models.Users{}
	err = config.DB.Model(&models.Users{}).Where("id = ?", claims.UserId).First(user).Error
	if err != nil || reflect.DeepEqual(user, &pb.Users{}) {
		log.Println("Failed to get users:", err.Error())
		return &pb.UserVerifyResponse{Valid: false, Message: "Invalid Credentials!"}, errors.New("Invalid credentials")
	}

	return &pb.UserVerifyResponse{Valid: true, Message: "Authenticated!", UserId: uint64(user.ID)}, nil
}

func generateToken(user *models.Users) (string, error) {
	now := time.Now()
	expiry := time.Now().Add(time.Hour * 24 * 2)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		UserId: uint64(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"go-gin-microservice-grpc"},
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseToken(tokenString, secret string) (claims AuthClaims, err error) {
	decodedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return AuthClaims{}, err
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok && decodedToken.Valid &&
		claims.VerifyAudience("go-gin-microservice-grpc", true) &&
		claims.VerifyExpiresAt(time.Now().Unix(), true) &&
		claims.VerifyIssuedAt(time.Now().Unix(), true) {

		authClaims := AuthClaims{}
		b, err := json.Marshal(claims)
		if err != nil {
			return AuthClaims{}, err
		}
		err = json.Unmarshal(b, &authClaims)
		if err != nil {
			return AuthClaims{}, err
		}
		return authClaims, nil
	}
	return AuthClaims{}, err
}
