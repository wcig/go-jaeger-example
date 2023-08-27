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

func Init(serviceName string) (opentracing.Tracer, io.Closer) {
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
	tracerCloser = closer
	return tracer, closer
}

func Close() {
	err := tracerCloser.Close()
	if err != nil {
		fmt.Println("jaeger close err:", err)
	}
}
