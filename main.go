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
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DBNAME")
	collectionName := os.Getenv("MONGODB_COLLECTION")
	port := os.Getenv("PORT")
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

	userRepo := repo.NewMongoUserRepository(client.Database(dbName).Collection(collectionName))
	userService := service.CreateUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.UserRoutes(router)

	authService := service.CreateAuthService(userService)
	authHandler := handler.NewAuthHandler(authService)
	authHandler.AuthRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Not Found",
		})
	})

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

	router.Run(":" + port)
}
