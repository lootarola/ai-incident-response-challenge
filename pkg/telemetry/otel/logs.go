package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func NewLoggerProvider(ctx context.Context, endpoint, serviceName, serviceVersion string) (*sdklog.LoggerProvider, error) {
	exp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create log exporter: %w", err)
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

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
		sdklog.WithResource(res),
	)
	return lp, nil
}

type Logger struct {
	logger log.Logger
}

func NewLogger(lp log.LoggerProvider, name string) *Logger {
	return &Logger{logger: lp.Logger(name)}
}

func (l *Logger) Info(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.emit(ctx, log.SeverityInfo, msg, attrs...)
}

func (l *Logger) Error(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.emit(ctx, log.SeverityError, msg, attrs...)
}

func (l *Logger) emit(ctx context.Context, severity log.Severity, msg string, attrs ...log.KeyValue) {
	var r log.Record
	r.SetTimestamp(time.Now())
	r.SetBody(log.StringValue(msg))
	r.SetSeverity(severity)

	all := make([]log.KeyValue, 0, len(attrs)+1)
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		all = append(all, log.String("trace_id", sc.TraceID().String()))
	}
	all = append(all, attrs...)
	r.AddAttributes(all...)

	l.logger.Emit(ctx, r)
}

func KV(k, v string) log.KeyValue {
	return log.String(k, v)
}

func KVFloat(k string, v float64) log.KeyValue {
	return log.Float64(k, v)
}

func KVInt(k string, v int) log.KeyValue {
	return log.Int(k, v)
}
