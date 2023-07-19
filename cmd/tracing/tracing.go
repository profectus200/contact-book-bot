package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func InitTracing(serviceName string, logger *zap.Logger) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}

	opentracing.SetGlobalTracer(tracer)
}
