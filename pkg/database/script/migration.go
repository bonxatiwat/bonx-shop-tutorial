package main

import (
	"context"
	"log"
	"os"

	"github.com/bonxatiwat/bonx-shop-tutorial/config"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/database/migration"
)

func main() {
	ctx := context.Background()
	_ = ctx

	// Initialize config
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path id required")
		}

		return os.Args[1]
	}())

	switch cfg.App.Name {
	case "player":
		migration.PlayerMigrate(ctx, &cfg)
	case "auth":
		migration.AuthMigrate(ctx, &cfg)
	case "item":
		migration.ItemMigrate(ctx, &cfg)
	case "inventory":
		migration.InventoryMigrate(ctx, &cfg)
	case "payment":
		migration.PaymentMigrate(ctx, &cfg)
	}
}
