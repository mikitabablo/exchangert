package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikitabablo/exchangert/config"
	"github.com/mikitabablo/exchangert/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config is not loaded: %v", err)
	}

	application := app.NewApp(cfg)
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = application.Stop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
