package main

import (
	"context"
	"log"
	"os/signal"
	"snappfood/cmd/command"
	"snappfood/internal/config"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	const description = "mabna stock market"
	root := &cobra.Command{Short: description}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	root.AddCommand(
		command.Server{}.Command(ctx, cfg),
	)

	if err := root.Execute(); err != nil {
		log.Fatalf("failed to execute root command: \n %s", err.Error())
	}
}
