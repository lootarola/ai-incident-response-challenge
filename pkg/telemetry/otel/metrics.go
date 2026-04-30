package otel

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewMeterProvider(ctx context.Context, endpoint, serviceName, serviceVersion string) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(15*time.Second))),
		sdkmetric.WithResource(res),
	)

	if err := registerRuntimeMetrics(ctx, mp); err != nil {
		return nil, fmt.Errorf("register runtime metrics: %w", err)
	}

	return mp, nil
}

func registerRuntimeMetrics(ctx context.Context, mp *sdkmetric.MeterProvider) error {
	meter := mp.Meter("runtime")

	_, err := meter.Int64ObservableGauge("go.goroutine.count",
		metric.WithDescription("Number of goroutines that currently exist."),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(runtime.NumGoroutine()))
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("register goroutine gauge: %w", err)
	}

	_, err = meter.Int64ObservableGauge("process.heap.inuse_bytes",
		metric.WithDescription("Bytes in in-use spans."),
		metric.WithUnit("By"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			o.Observe(int64(ms.HeapInuse))
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("register heap gauge: %w", err)
	}

	return nil
}
