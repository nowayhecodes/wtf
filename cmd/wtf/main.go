package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nowayhecodes/wtf/internal/app"
	"github.com/nowayhecodes/wtf/internal/config"
	"github.com/nowayhecodes/wtf/internal/correction"
	"github.com/nowayhecodes/wtf/internal/history"
	"github.com/nowayhecodes/wtf/internal/shell"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	histParser := history.NewParser()
	detector := correction.NewDetector(cfg)
	executor := shell.NewExecutor()

	application := app.New(cfg, histParser, detector, executor)
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
