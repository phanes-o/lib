package trace

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"

	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

var defaultProviderOptions = &ProviderOptions{
	name:         serviceName,
	env:          "dev",
	url:          collectorURL,
	insecure:     len(insecure) > 0,
	otelProtocol: OtelGRPC,
}

func NewTraceProvider(providerType ProviderType, opts ...ProviderOption) (*trace.TracerProvider, error)  {
	switch providerType {
	case ProviderTypeJaeger:
		return jaegerProvider(opts...)
	case ProviderTypeZipkin:
		return zipkinProvider(opts...)
	case ProviderTypeOTLP:
		return otelProvider(opts...)
	default:
		return nil, errors.New("invalid provider type")
	}
}

func jaegerProvider(opts ...ProviderOption) (*trace.TracerProvider, error) {
	config := defaultProviderOptions
	for _, fn := range opts {
		fn.Apply(config)
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.url)))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.name),
			attribute.String("environment", config.env),
			attribute.String("version", config.version),
		)),
	)
	return tp, nil
}

func zipkinProvider(opts ...ProviderOption)  (*trace.TracerProvider, error) {
	config := defaultProviderOptions
	for _, fn := range opts {
		fn.Apply(config)
	}

	exporter, err := zipkin.New(config.url)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.name),
			attribute.String("environment", config.env),
			attribute.String("version", config.version),
		)),
	)
	return tp, nil
}

func otelProvider(opts ...ProviderOption) (*trace.TracerProvider, error) {
	config := defaultProviderOptions
	for _, fn := range opts {
		fn.Apply(config)
	}


	var client otlptrace.Client
	switch config.otelProtocol {
	case OtelGRPC:
		secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if config.insecure {
			secureOption = otlptracegrpc.WithInsecure()
		}
		client = otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.url),
		)
	case OtelHTTP:
		client = otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(config.url),
		)
	}

	exporter, err := otlptrace.New(
		context.Background(),
		client,
	)

	if err != nil {
		return nil, err
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", config.name),
			attribute.String("library.language", "go"),
			attribute.String("service.version", config.version),
		),
	)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)

	return tp, nil
}
