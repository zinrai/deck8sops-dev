package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/operations"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

func main() {
	configFile := flag.String("config", "operations.yaml", "Config file path")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	flag.CommandLine.Parse(os.Args[2:])

	logger := utils.NewLogger(*verbose)

	logger.Info("Reading configuration from %s", *configFile)
	cfg, err := config.ReadFromFile(*configFile)
	if err != nil {
		logger.Error("Failed to read config: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	switch cmd {
	case "create":
		logger.Info("Starting operation creation")
		if err := operations.Create(ctx, cfg, logger); err != nil {
			logger.Error("Error creating operations: %v", err)
			os.Exit(1)
		}
		logger.Info("All operations created successfully")

	case "delete":
		logger.Info("Starting operation deletion")
		if err := operations.Delete(ctx, cfg, logger); err != nil {
			logger.Error("Error deleting operations: %v", err)
			os.Exit(1)
		}
		logger.Info("All operations deleted successfully")

	default:
		logger.Error("Unknown command: %s", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("deck8sops-dev - Declarative Kubernetes Operator setup for Kind clusters")
	fmt.Println("\nUsage:")
	fmt.Println("  deck8sops-dev [command] [flags]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  create      Create operations defined in config file")
	fmt.Println("  delete      Delete operations defined in config file")
	fmt.Println("\nFlags:")
	fmt.Println("  -config string   Config file path (default \"operations.yaml\")")
	fmt.Println("  -verbose         Enable verbose output")
	fmt.Println("\nExample:")
	fmt.Println("  deck8sops-dev create -config=my-operations.yaml")
}
