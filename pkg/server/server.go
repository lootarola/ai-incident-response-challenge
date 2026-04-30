package server

import (
	"context"
	"fmt"
	"os"
	"sync"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server interface {
	Start(ctx context.Context) error
}

// Handler is a backend-agnostic marker. The chosen backend type-asserts to its
// own routing shape (e.g. RegisterRoutes(gin.IRouter) for the gin backend).
type Handler interface{}

type Config struct {
	Handlers       []Handler
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
}

type Factory func(cfg Config) (Server, error)

var (
	mu        sync.RWMutex
	factories = map[string]Factory{}
)

func Register(name string, f Factory) {
	mu.Lock()
	defer mu.Unlock()
	factories[name] = f
}

func New(cfg Config) (Server, error) {
	backend := os.Getenv("SERVER_BACKEND")
	mu.RLock()
	f, ok := factories[backend]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("SERVER_BACKEND %q not registered", backend)
	}
	return f(cfg)
}
