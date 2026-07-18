// Package main is the single entry point.
//
// @title           Go Base Template API
// @version         1.0
// @description     Horizontal-layer modular monolith base template (no repository pattern, no outbox, no cmd/worker).
// @BasePath        /api/v1
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/furkancmn57/go-base-template/src/common/apperr"
	"github.com/furkancmn57/go-base-template/src/config"
	"github.com/furkancmn57/go-base-template/src/constants"
	todoconsumer "github.com/furkancmn57/go-base-template/src/consumers/todo"
	v1 "github.com/furkancmn57/go-base-template/src/controllers/v1"
	"github.com/furkancmn57/go-base-template/src/extensions"
	todoservice "github.com/furkancmn57/go-base-template/src/services/todo"

	_ "github.com/furkancmn57/go-base-template/src/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("main: failed to load config: %v", err)
	}

	gormDB, err := extensions.AddDatabase(cfg.Postgres)
	if err != nil {
		log.Fatalf("main: database: %v", err)
	}

	redisClient, err := extensions.AddRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("main: redis: %v", err)
	}
	defer redisClient.Close()

	mq, err := extensions.AddRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("main: rabbitmq: %v", err)
	}
	defer mq.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return apperr.WriteHTTP(c, apperr.Internal(err))
		},
	})

	extensions.RegisterHealth(app, gormDB, redisClient)
	extensions.RegisterOpenAPI(app)

	todoService := todoservice.NewService(gormDB, mq)
	api := app.Group("/api/" + constants.APIVersion)
	v1.NewTodoController(todoService).Register(api)

	if err := todoconsumer.NewLogWhenTodoCompleted().Register(mq); err != nil {
		log.Fatalf("main: failed to register todo consumer: %v", err)
	}

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Printf("main: server stopped: %v", err)
		}
	}()
	log.Printf("main: listening on port %s (env=%s)", cfg.AppPort, cfg.AppEnv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("main: shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("main: error during shutdown: %v", err)
	}
}
