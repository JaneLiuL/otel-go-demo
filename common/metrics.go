package common

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metrics"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.29.0"
)

func InitMetrics(serviceName string) (metrics.MeterProvider, error) {
	exporter, err := prometheus.New(prometheus.WithoutUnits())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metrics.NewMeterProvider(
		metrics.WithReader(exporter),
		metrics.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exporter)
		server := &http.Server{
			Addr:    ":8888",
			Handler: mux,
		}
		log.Println("Starting metrics server on :8888")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed
		{
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()
	log.Println("Metrics initialized successfully")
	return meterProvider, nil
}