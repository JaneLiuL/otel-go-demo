package main

import (
    "log"
    "net/http"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "github.com/JaneLiuL/otel-go-demo/common"
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "context"
)
func main() {
    tp, err := common.InitTracer("otel-go-demo", "http://localhost:14268/api/traces")
    if err != nil {
        log.Fatalf("failed to initialize tracer: %v", err)
    }
    defer tp.Shutdown(context.Background())


    mp, err := common.InitMetrics("otel-go-demo")
    if err != nil {
        log.Fatalf("failed to initialize metrics: %v", err)
    }
    defer mp.Shutdown(context.Background())
    log.Println("OpenTelemetry initialized successfully")
    meter := mp.Meter("otel-go-demo")
    counter, err := meter.Int64Counter("requests_total",
        common.WithDescription("Total number of requests received"),
    )
    if err != nil {
        log.Fatalf("failed to create counter: %v", err)
    }
    counter.Add(context.Background(), 1, resource.WithAttributes(
        semconv.ServiceNameKey.String("otel-go-demo"),
        semconv.ServiceVersionKey.String("1.0.0"),
        semconv.DeploymentEnvironmentKey.String("production"),
    ))

    mux := http.NewServeMux()
    mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, OpenTelemetry!"))
    }), "HelloHandler"))
    server := &http.Server{
        Addr:    ":8080",
        Handler: otelhttp.NewHandler(mux, "HTTPServer"), 
    }
    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}