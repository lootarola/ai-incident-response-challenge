package otel

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	LoggerProvider *sdklog.LoggerProvider
}

type Scope struct {
	Logger *Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}

func (t *Telemetry) For(name string) *Scope {
	return &Scope{
		Logger: NewLogger(t.LoggerProvider, name),
		Tracer: t.TracerProvider.Tracer(name),
		Meter:  t.MeterProvider.Meter(name),
	}
}

func New(ctx context.Context) (*Telemetry, func(), error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	serviceVersion := os.Getenv("SERVICE_VERSION")

	tp, err := NewTracerProvider(ctx, endpoint, serviceName, serviceVersion)
	if err != nil {
		return nil, nil, err
	}

	mp, err := NewMeterProvider(ctx, endpoint, serviceName, serviceVersion)
	if err != nil {
		return nil, nil, err
	}

	lp, err := NewLoggerProvider(ctx, endpoint, serviceName, serviceVersion)
	if err != nil {
		return nil, nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	shutdown := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown tracer: %v", err)
		}
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown meter: %v", err)
		}
		if err := lp.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown logger: %v", err)
		}
	}

	return &Telemetry{
		TracerProvider: tp,
		MeterProvider:  mp,
		LoggerProvider: lp,
	}, shutdown, nil
}
