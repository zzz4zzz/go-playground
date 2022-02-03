package trace

import (
	"context"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type JaegerConfiguration struct {
	Endpoint    string
	ServiceName string
	Enabled     bool
}

type Provider struct {
	Provider    oteltrace.TracerProvider
	Propagators propagation.TextMapPropagator
}

func NewProvider(config JaegerConfiguration) (Provider, error) {
	if !config.Enabled {
		return Provider{Provider: oteltrace.NewNoopTracerProvider()}, nil
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Endpoint)),
	)
	if err != nil {
		return Provider{}, err
	}

	prv := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
		)),
	)
	propagators := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return Provider{Provider: prv, Propagators: propagators}, nil
}
func (p Provider) Close(ctx context.Context) error {
	if prv, ok := p.Provider.(*sdktrace.TracerProvider); ok {
		return prv.Shutdown(ctx)
	}

	return nil
}
