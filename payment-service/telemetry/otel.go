package telemetry

import (
    "context"
    "errors"
    "fmt"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/log/global"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdklog "go.opentelemetry.io/otel/sdk/log"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Setup(ctx context.Context, serviceName, collectorURL string) (shutdown func(context.Context) error, err error) {
    res, err := resource.New(ctx,
        resource.WithAttributes(semconv.ServiceName(serviceName)),
    )
    if err != nil {
        return nil, fmt.Errorf("resource: %w", err)
    }

    opts := []otlptracehttp.Option{
        otlptracehttp.WithEndpoint(collectorURL),
        otlptracehttp.WithInsecure(),
    }

    // --- Traces ---
    traceExp, err := otlptracehttp.New(ctx, opts...)
    if err != nil {
        return nil, fmt.Errorf("trace exporter: %w", err)
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(traceExp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})

    // --- Metrics ---
    metricExp, err := otlpmetrichttp.New(ctx,
        otlpmetrichttp.WithEndpoint(collectorURL),
        otlpmetrichttp.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("metric exporter: %w", err)
    }
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp,
            sdkmetric.WithInterval(10*time.Second),
        )),
        sdkmetric.WithResource(res),
    )
    otel.SetMeterProvider(mp)

    // --- Logs ---
    logExp, err := otlploghttp.New(ctx,
        otlploghttp.WithEndpoint(collectorURL),
        otlploghttp.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("log exporter: %w", err)
    }
    lp := sdklog.NewLoggerProvider(
        sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
        sdklog.WithResource(res),
    )
    global.SetLoggerProvider(lp)

    shutdown = func(ctx context.Context) error {
        return errors.Join(
            tp.Shutdown(ctx),
            mp.Shutdown(ctx),
            lp.Shutdown(ctx),
        )
    }
    return shutdown, nil
}
