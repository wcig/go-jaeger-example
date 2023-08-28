package jaeger

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	SamplerConstTrue = 1
	Endpoint         = "http://127.0.0.1:14268/api/traces"
)

var (
	tracerCloser io.Closer
)

func Init(serviceName string) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: SamplerConstTrue,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: Endpoint,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	tracerCloser = closer
}

func Close() {
	if err := tracerCloser.Close(); err != nil {
		fmt.Println("jaeger close err:", err)
	}
}
