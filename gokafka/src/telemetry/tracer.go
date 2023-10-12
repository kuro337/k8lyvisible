package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupTracer(ctx context.Context, logger *slog.Logger, svc string, url string) (cleanup func(context.Context), err error) {
	logger.Info("Creating gRPC OTLP Trace Exporter")

	shutdown, _, err := CreateOTLPTraceExporterGRPC(ctx, svc, url) // svc, url "go-app","0.0.0.0:4317"
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Return this cleanup function for caller to defer

	return func(ctx context.Context) {
		shutdown(ctx) // Execute the shutdown without a context
		logger.Info("Tracer Shutdown complete")
	}, nil
}

func CreateOTLPTraceExporterGRPC(ctx context.Context, serviceName string, grpcEndpoint string) (cleanup func(ctx context.Context), cntx context.Context, err error) {
	//	ctx := context.Background()

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
	conn, err := grpc.DialContext(ctx, grpcEndpoint, // "0.0.0.0:4317"
		// TODO: Add TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Setup Trace Exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register Trace Exporter with a TraceProvider - using Batch Span Processing
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown Flushes remaining spans and Shuts Down the Exporter
	cleanup = func(ctx context.Context) {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			//	l.Fatal(err)
		}
	}

	return cleanup, ctx, nil
}
