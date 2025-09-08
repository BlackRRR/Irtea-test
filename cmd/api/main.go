package main

import (
	"context"
	"log"
	"github.com/BlackRRR/Irtea-test/internal/app"
	"github.com/joho/godotenv"
)

func init() {
	app.InternalInit()
}

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	ctx := context.Background()

	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// initialize application
	a := app.Init(ctx, cfg)

	// run application
	a.Run(ctx)
}
