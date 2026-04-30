package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lootarola/ai-incident-response-challenge/pkg/server"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func init() {
	server.Register("gin", newServer)
}

type handler interface {
	RegisterRoutes(r gin.IRouter)
}

type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
	port       string
}

func newServer(cfg server.Config) (server.Server, error) {
	handlers := make([]handler, 0, len(cfg.Handlers))
	for _, h := range cfg.Handlers {
		gh, ok := h.(handler)
		if !ok {
			return nil, fmt.Errorf("handler %T does not implement gin handler", h)
		}
		handlers = append(handlers, gh)
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(otelgin.Middleware("incident-response-api",
		otelgin.WithTracerProvider(cfg.TracerProvider),
	))
	engine.Use(metricsMiddleware(cfg.MeterProvider.Meter("http.server")))

	engine.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, h := range handlers {
		h.RegisterRoutes(engine)
	}

	port := os.Getenv("HTTP_PORT")
	return &Server{
		engine: engine,
		port:   port,
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: engine,
		},
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutCtx)
	}
}

func metricsMiddleware(meter metric.Meter) gin.HandlerFunc {
	reqDuration, _ := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Duration of HTTP server requests."),
	)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start).Seconds()

		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.Int("http.response.status_code", c.Writer.Status()),
		}
		reqDuration.Record(c.Request.Context(), dur, metric.WithAttributes(attrs...))
	}
}
