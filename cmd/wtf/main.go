package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nowayhecodes/wtf/internal/app"
	"github.com/nowayhecodes/wtf/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	application := app.New(cfg)
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
