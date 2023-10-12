package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupMetrics(ctx context.Context, logger *slog.Logger, svc string, url string) (cleanup func(context.Context), err error) {
	logger.Info("Creating gRPC OTLP Metric Exporter")

	shutdown, _, err := CreateOTLPMetricExporterGRPC(ctx, logger, svc, url) // svc, url "go-app","0.0.0.0:4318"
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Return this cleanup function for the caller to defer

	return func(ctx context.Context) {
		shutdown(ctx)
		logger.Info("Metric Exporter Shutdown complete")
	}, nil
}

func CreateOTLPMetricExporterGRPC(ctx context.Context, logger *slog.Logger, serviceName string, grpcEndpoint string) (cleanup func(ctx context.Context), cntx context.Context, err error) {
	// Create a resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName), // service name
		),
	)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Create a gRPC connection to the collector
	conn, err := grpc.DialContext(ctx, grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Create an OTLP metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Create a metric meter provider and configure it to use the exporter
	meterProvider := sdk.NewMeterProvider(
		sdk.WithResource(res),
		sdk.WithReader(sdk.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider)

	// Shutdown function for cleanup
	cleanup = func(ctx context.Context) {
		if err := meterProvider.Shutdown(ctx); err != nil {
			logger.Error(err.Error())
		}
	}

	return cleanup, ctx, nil
}

func CaptureHeapAllocations(meter metric.Meter) {
	if _, err := meter.Int64ObservableGauge(
		"memory.heap",
		metric.WithDescription(
			"Memory usage of the allocated heap objects.",
		),
		metric.WithUnit("By"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			o.Observe(int64(m.HeapAlloc))
			return nil
		}),
	); err != nil {
		panic(err)
	}
}
