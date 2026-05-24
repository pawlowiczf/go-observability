package telemetry

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	reserveRequestCounter    metric.Int64Counter
	reserveErrorCounter      metric.Int64Counter
	reserveDurationHistogram metric.Float64Histogram
	processCPUPercentGauge   metric.Float64ObservableGauge
	currentProcess           *process.Process
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

	processCPUPercentGauge, err = meter.Float64ObservableGauge(
		"process_cpu_percent",
		metric.WithDescription("CPU usage percentage of the process"),
	)
	if err != nil {
		return err
	}

	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return err
	}
	currentProcess = p

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		if currentProcess != nil {
			if cpuPercent, err := currentProcess.Percent(100 * time.Millisecond); err == nil {
				observer.ObserveFloat64(processCPUPercentGauge, cpuPercent)
			}
		}
		return nil
	}, processCPUPercentGauge)
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
