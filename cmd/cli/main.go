package main

import (
	"flag"
	"fmt"

	"context"
	"time"

	"github.com/mises-id/mises-swapsvc/app/services/swap"
	"github.com/mises-id/mises-swapsvc/lib/db"

	// This Service
	"github.com/mises-id/mises-swapsvc/handlers"
	"github.com/mises-id/mises-swapsvc/svc/server"
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	fmt.Println("setup mongo...")
	db.SetupMongo(ctx)
	//models.EnsureIndex()
	swap.SyncSwapOrder(ctx)
	cfg := server.DefaultConfig
	cfg = handlers.SetConfig(cfg)

	server.Run(cfg)
}
