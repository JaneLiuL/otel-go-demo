package common

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(serviceName string, jagerEndpoint string) (*trace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jagerEndpoint)))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
