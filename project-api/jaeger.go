package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	TraceHeader = "Trace-ID"
	SpanKey     = "gin-span"
	Component   = "GIN"
)

func InitJaeger() {
	cfg := &config.Configuration{
		ServiceName: "project-api",
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
		},
	}
	tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
}

func JaegerTrace(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	fullPath := c.FullPath()
	traceId := c.GetHeader(TraceHeader)

	var span opentracing.Span
	if traceId == "" {
		span = tracer.StartSpan(fullPath)
	} else {
		upstreamCtx, err := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		span = opentracing.StartSpan(fullPath, ext.RPCServerOption(upstreamCtx))
	}
	defer span.Finish()

	ext.HTTPMethod.Set(span, c.Request.Method)
	ext.HTTPUrl.Set(span, fullPath)
	ext.Component.Set(span, Component)
	c.Set(SpanKey, opentracing.ContextWithSpan(c, span))
	c.Next()
	ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
}
