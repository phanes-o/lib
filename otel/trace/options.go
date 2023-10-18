package trace

type Protocol string
type ProviderType string

const (
	ProviderTypeJaeger ProviderType = "jaeger"
	ProviderTypeZipkin ProviderType = "zipkin"
	ProviderTypeOTLP   ProviderType = "otlp"
)

const (
	OtelGRPC Protocol = "grpc"
	OtelHTTP Protocol = "http"
)

type ProviderOptions struct {
	name         string
	env          string
	version      string
	url          string
	insecure     bool
	otelProtocol Protocol // grpc or http
}

type ProviderOption interface {
	Apply(o *ProviderOptions)
}

type providerOptionFunc func(o *ProviderOptions)

func (f providerOptionFunc) Apply(o *ProviderOptions) {
	f(o)
}

func WithName(name string) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.name = name
	})
}

func WithEnv(env string) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.env = env
	})
}

func WithURL(url string) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.url = url
	})
}

func WithInsecure(insecure bool) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.insecure = insecure
	})
}

func WithVersion(version string) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.version = version
	})
}

func WithOtelProtocol(protocol Protocol) ProviderOption {
	return providerOptionFunc(func(o *ProviderOptions) {
		o.otelProtocol = protocol
	})
}
