package telemetry

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	reserveRequestCounter    metric.Int64Counter
	reserveErrorCounter      metric.Int64Counter
	reserveDurationHistogram metric.Float64Histogram
	processMemoryGauge       metric.Int64ObservableGauge
	processorGoroutineGauge  metric.Int64ObservableGauge
	metricsInitialized       bool
)

func InitMetrics(serviceName string) error {
	meter := otel.Meter(serviceName)
	var err error

	reserveRequestCounter, err = meter.Int64Counter(
		"inventory_reserve_requests_total",
		metric.WithDescription("Number of inventory reserve requests"),
	)
	if err != nil {
		return err
	}

	reserveErrorCounter, err = meter.Int64Counter(
		"inventory_reserve_errors_total",
		metric.WithDescription("Number of failed inventory reserve requests"),
	)
	if err != nil {
		return err
	}

	reserveDurationHistogram, err = meter.Float64Histogram(
		"inventory_reserve_duration_seconds",
		metric.WithDescription("Duration of inventory reserve handler in seconds"),
	)
	if err != nil {
		return err
	}

	processMemoryGauge, err = meter.Int64ObservableGauge(
		"process_memory_bytes",
		metric.WithDescription("Process memory usage in bytes"),
	)
	if err != nil {
		return err
	}

	processorGoroutineGauge, err = meter.Int64ObservableGauge(
		"process_goroutines",
		metric.WithDescription("Number of goroutines"),
	)
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		observer.ObserveInt64(processMemoryGauge, int64(mem.Alloc))
		observer.ObserveInt64(processorGoroutineGauge, int64(runtime.NumGoroutine()))
		return nil
	}, processMemoryGauge, processorGoroutineGauge)
	if err != nil {
		return err
	}

	metricsInitialized = true
	return nil
}

func RecordReserveRequest(ctx context.Context, duration time.Duration, success bool) {
	if !metricsInitialized {
		return
	}

	status := "ok"
	if !success {
		status = "error"
	}

	attrs := []attribute.KeyValue{
		attribute.String("http.method", http.MethodPost),
		attribute.String("http.target", "/reserve"),
		attribute.String("status", status),
	}

	reserveRequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	reserveDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
	if !success {
		reserveErrorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}
