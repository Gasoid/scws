package tracing

import (
	"io"
	"log"

	opentracing "github.com/opentracing/opentracing-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

const (
	serviceName = "scws"
)

func JaegerInit() (io.Closer, error) {
	metricsFactory := prometheus.New()
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		ServiceName: serviceName,
	}
	jLogger := jaegerlog.StdLogger
	_, err := cfg.FromEnv()
	if err != nil {
		log.Println("tracing.JaegerInit", err.Error())
		return nil, err
	}
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(metricsFactory),
	)
	if err != nil {
		return nil, err
	}
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}
