package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiber-boilerplate/internal/app/config"
	"fiber-boilerplate/internal/app/handlers"
	"fiber-boilerplate/internal/app/middleware"
	"fiber-boilerplate/internal/models"
	"fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/setting"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Start() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using system environment variables")
	}

	config.Setup()
	models.Setup()

	f := fiber.New()

	middleware.Register(f)
	handlers.Register(f)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", setting.Runtime.Port)
		err := f.Listen(addr)
		switch {
		case err == nil, errors.Is(err, http.ErrServerClosed):
			// normal
		default:
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
	logging.Info("Shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Server.GracefulTimeout)*time.Second)
	defer cancel()

	if err := f.ShutdownWithContext(ctx); err != nil {
		logging.Error(err, "Failed to shutdown completely")
	}
}
