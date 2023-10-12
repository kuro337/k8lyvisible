package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"kafkago/simulation"
	"kafkago/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	collectorEndpoint = "otel-collector-cluster-opentelemetry-collector.default.svc.cluster.local:4317"
	serviceName       = "kafkago"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Hello World")

	ctx := context.Background()

	// Tracing Telemetry
	tracercleanup, err := telemetry.SetupTracer(ctx, logger, serviceName, collectorEndpoint)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer tracercleanup(ctx)

	tracer := otel.GetTracerProvider().Tracer(serviceName)
	_, span := tracer.Start(ctx, "BINARY_EXECUTED")

	defer func() {
		span.AddEvent("EXIT", trace.WithAttributes(
			attribute.String("BINARY_EXIT_STATUS", "SUCCESS"),
			attribute.String("BINARY_EXIT_TIMESTAMP", time.Now().String()),
		))
		span.End()
		logger.Info("Span Ended")
	}()

	// Metrics Telemetry
	metriccleanup, err := telemetry.SetupMetrics(ctx, logger, serviceName, collectorEndpoint)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer metriccleanup(ctx)

	meter := otel.Meter(serviceName)
	telemetry.CaptureHeapAllocations(meter)

	simulation.SimulateOperations()

	logger.Info("App Exited Successfully")
}

/*
brokers: ["redpanda-0.redpanda.redpanda.svc.cluster.local.:9093"],

OTEL Collector
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
value: "http://otel-collector-cluster-opentelemetry-collector:4317"

*/
