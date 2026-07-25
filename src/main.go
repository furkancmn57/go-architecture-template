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
	"github.com/furkancmn57/go-architecture-template/src/extensions"
	appgraphql "github.com/furkancmn57/go-architecture-template/src/graphql"
	todoservice "github.com/furkancmn57/go-architecture-template/src/services/todo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	logger.Init(cfg.AppEnv)

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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	todoService := todoservice.NewService()
	schema, err := appgraphql.NewSchema(todoService)
	if err != nil {
		slog.Error("graphql schema", "error", err)
		os.Exit(1)
	}
	extensions.RegisterGraphQL(app, schema)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			slog.Error("server stopped", "error", err)
		}
	}()
	slog.Info("listening", "port", cfg.AppPort, "env", cfg.AppEnv, "graphql", "/graphql")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
