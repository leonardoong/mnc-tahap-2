package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/leonardoong/e-wallet/config"
	"github.com/leonardoong/e-wallet/internal/publisher"
	"github.com/leonardoong/e-wallet/internal/queue"
	"github.com/leonardoong/e-wallet/internal/repository"
	"github.com/leonardoong/e-wallet/internal/routes"
	"github.com/leonardoong/e-wallet/internal/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables instead")
	}

	cfg := config.LoadConfig()

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	dbConn, err := sql.Open(`mysql`, connection)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	var cache *redis.Pool
	cache = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				return nil, fmt.Errorf("ERROR connect redis | %v", err)
			}
			if cfg.RedisPassword != "" {
				if _, err := c.Do("AUTH", cfg.RedisPassword); err != nil {
					return nil, fmt.Errorf("ERROR connect redis | %v", err)
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
	cacheConn, err := cache.Dial()
	if err != nil {
		log.Fatal("got error when connect redis", err)
		return
	}
	cfg.CachePool = cache
	defer cacheConn.Close()

	redisPublisher := publisher.NewPublisher("ewallet", cache)
	redisPublisher.Initialize()

	userRepo := repository.NewUserRepository(dbConn)
	walletRepo := repository.NewWalletRepository(dbConn)
	transactionRepo := repository.NewTransactionRepository(dbConn, redisPublisher)

	userService := service.NewAuthService(cfg, userRepo)
	transactionService := service.NewTransactionService(cfg, dbConn, transactionRepo, walletRepo, userRepo)

	redisConsumer := queue.NewQueue(cfg, transactionService)
	redisConsumer.Initialize()

	router := gin.Default()

	routes.SetupRoutes(router, userService, transactionService)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
