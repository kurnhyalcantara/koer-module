package tracing

import (
	"context"
	"fmt"

	"github.com/koer/koer-module/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	tp *sdktrace.TracerProvider
}

func NewProvider(ctx context.Context, cfg config.TracingConfig) (*Provider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	var exporter sdktrace.SpanExporter
	if cfg.Enabled && cfg.Endpoint != "" {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("creating OTLP exporter: %w", err)
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("creating stdout exporter: %w", err)
		}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return &Provider{tp: tp}, nil
}

func (p *Provider) Tracer(name string) trace.Tracer {
	return p.tp.Tracer(name)
}

func (p *Provider) Shutdown(ctx context.Context) error {
	if err := p.tp.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down tracer provider: %w", err)
	}
	return nil
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
