package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	handler "7-solutions-test-golang/adapters/handlers"
	repo "7-solutions-test-golang/adapters/repositories/mongo"
	service "7-solutions-test-golang/core/service"
	middleware "7-solutions-test-golang/middlewares"
	util "7-solutions-test-golang/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DBNAME")
	collectionName := os.Getenv("MONGODB_COLLECTION")
	delay, err := strconv.Atoi(os.Getenv("DELAY_SECOND"))
	if err != nil {
		panic(err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	router := gin.Default()
	util.RegisterValidation()

	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	indexHandler := handler.NewIndexHandler()
	indexHandler.IndexHandler(router)

	userRepo := repo.NewMongoUserRepository(client.Database(dbName).Collection(collectionName))
	userService := service.CreateUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.UserRoutes(router)

	authService := service.CreateAuthService(userService)
	authHandler := handler.NewAuthHandler(authService)
	authHandler.AuthRoutes(router)

	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				amount, err := userService.GetCount()
				if err != nil {
					amount = 0
				}
				fmt.Printf("Amount Users : %d\n", amount)
			}
		}
	}()

	defer func() {
		ticker.Stop()
		done <- true
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run("0.0.0.0:" + port)
}
