package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type contextKey int

const (
	domainKey contextKey = iota
	operationKey
)

func WithDomain(ctx context.Context, domain string) context.Context {
	return context.WithValue(ctx, domainKey, domain)
}

func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func Domain(ctx context.Context) string {
	v, _ := ctx.Value(domainKey).(string)
	return v
}

func Operation(ctx context.Context) string {
	v, _ := ctx.Value(operationKey).(string)
	return v
}

func TraceID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		return sc.TraceID().String()
	}
	return ""
}
