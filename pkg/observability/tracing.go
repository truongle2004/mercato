// Package observability configures process-wide telemetry exporters.
package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/truongle2004/mercato/pkg/config"
)

// SetupTracing configures OTLP trace export and returns a shutdown function.
// When telemetry is disabled, the returned function is a no-op.
func SetupTracing(ctx context.Context, cfg *config.Schema) (func(context.Context) error, error) {
	if !cfg.TelemetryEnabled {
		return func(context.Context) error { return nil }, nil
	}

	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.TelemetryOTLPEndpoint),
	}
	if cfg.TelemetryOTLPInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceName(cfg.TelemetryServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("create telemetry resource: %w", err)
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}
