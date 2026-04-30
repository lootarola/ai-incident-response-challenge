package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lootarola/ai-incident-response-challenge/internal/catalog"
	"github.com/lootarola/ai-incident-response-challenge/internal/order"
	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/server"
	"github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tel, shutdown, err := otel.New(ctx)
	if err != nil {
		log.Fatalf("telemetry setup: %v", err)
	}
	defer shutdown()

	orderRepo, closeOrderRepo, err := database.NewOrder(ctx)
	if err != nil {
		log.Fatalf("order database: %v", err)
	}
	defer func() {
		if err := closeOrderRepo(context.Background()); err != nil {
			log.Printf("close order db: %v", err)
		}
	}()

	catalogRepo, closeCatalogRepo, err := database.NewCatalog(ctx)
	if err != nil {
		log.Fatalf("catalog database: %v", err)
	}
	defer func() {
		if err := closeCatalogRepo(context.Background()); err != nil {
			log.Printf("close catalog db: %v", err)
		}
	}()

	order, err := order.New(orderRepo, tel.For("order"))
	if err != nil {
		log.Fatalf("order module: %v", err)
	}

	catalog, err := catalog.New(catalogRepo, tel.For("catalog"))
	if err != nil {
		log.Fatalf("catalog module: %v", err)
	}

	srv, err := server.New(server.Config{
		Handlers:       []server.Handler{order.Handler, catalog.Handler},
		TracerProvider: tel.TracerProvider,
		MeterProvider:  tel.MeterProvider,
	})
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	log.Printf("starting server on :%s", os.Getenv("HTTP_PORT"))
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("server: %v", err)
	}
}
