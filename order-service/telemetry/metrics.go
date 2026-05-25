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
	checkoutRequestCounter    metric.Int64Counter
	checkoutErrorCounter      metric.Int64Counter
	checkoutDurationHistogram metric.Float64Histogram
	processCPUPercentGauge    metric.Float64ObservableGauge
	processMemoryRSSGauge     metric.Int64ObservableGauge
	currentProcess            *process.Process
	metricsInitialized        bool
)

func InitMetrics(serviceName string) error {
	meter := otel.Meter(serviceName)
	var err error

	checkoutRequestCounter, err = meter.Int64Counter(
		"order_checkout_requests_total",
		metric.WithDescription("Number of order checkout requests"),
	)
	if err != nil {
		return err
	}

	checkoutErrorCounter, err = meter.Int64Counter(
		"order_checkout_errors_total",
		metric.WithDescription("Number of failed order checkout requests"),
	)
	if err != nil {
		return err
	}

	checkoutDurationHistogram, err = meter.Float64Histogram(
		"order_checkout_duration_seconds",
		metric.WithDescription("Duration of order checkout handler in seconds"),
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

	processMemoryRSSGauge, err = meter.Int64ObservableGauge(
		"process_memory_rss_bytes",
		metric.WithDescription("Resident set size of the process in bytes"),
		metric.WithUnit("By"),
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
		if currentProcess == nil {
			return nil
		}
		if cpuPercent, err := currentProcess.Percent(100 * time.Millisecond); err == nil {
			observer.ObserveFloat64(processCPUPercentGauge, cpuPercent)
		}
		if mem, err := currentProcess.MemoryInfo(); err == nil {
			observer.ObserveInt64(processMemoryRSSGauge, int64(mem.RSS))
		}
		return nil
	}, processCPUPercentGauge, processMemoryRSSGauge)
	if err != nil {
		return err
	}

	metricsInitialized = true
	return nil
}

func RecordCheckoutRequest(ctx context.Context, duration time.Duration, success bool) {
	if !metricsInitialized {
		return
	}

	status := "ok"
	if !success {
		status = "error"
	}

	attrs := []attribute.KeyValue{
		attribute.String("http.method", http.MethodPost),
		attribute.String("http.target", "/checkout"),
		attribute.String("status", status),
	}

	checkoutRequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	checkoutDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	if !success {
		checkoutErrorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}
