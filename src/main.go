// Package main is the single entry point.
//
// @title           Go Base Template API
// @version         1.0
// @description     Horizontal-layer modular monolith base template (no repository pattern, no outbox, no cmd/worker).
// @BasePath        /api/v1
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/furkancmn57/go-architecture-template/src/common/apperr"
	"github.com/furkancmn57/go-architecture-template/src/common/logger"
	"github.com/furkancmn57/go-architecture-template/src/config"
	"github.com/furkancmn57/go-architecture-template/src/constants"
	v1 "github.com/furkancmn57/go-architecture-template/src/controllers/v1"
	"github.com/furkancmn57/go-architecture-template/src/extensions"
	todoservice "github.com/furkancmn57/go-architecture-template/src/services/todo"

	_ "github.com/furkancmn57/go-architecture-template/src/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	logger.Init(cfg.AppEnv)

	gormDB, err := extensions.AddDatabase(cfg.Postgres)
	if err != nil {
		slog.Error("database", "error", err)
		os.Exit(1)
	}

	redisClient, err := extensions.AddRedis(cfg.Redis)
	if err != nil {
		slog.Error("redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return apperr.WriteHTTP(c, apperr.Internal(err))
		},
	})
	if logger.IsDev(cfg.AppEnv) {
		app.Use(fiberlogger.New(fiberlogger.Config{
			Format: "[${time}] ${status} ${method} ${path} ${latency}\n",
			Output: os.Stdout,
		}))
	}

	extensions.RegisterHealth(app, gormDB, redisClient)
	extensions.RegisterOpenAPI(app)

	todoService := todoservice.NewService(gormDB)
	api := app.Group("/api/" + constants.APIVersion)
	v1.NewTodoController(todoService).Register(api)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			slog.Error("server stopped", "error", err)
		}
	}()
	slog.Info("listening", "port", cfg.AppPort, "env", cfg.AppEnv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("error during shutdown", "error", err)
	}
}
