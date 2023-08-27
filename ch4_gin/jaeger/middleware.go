package jaeger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	SpanCtxKey = "SpanCtxKey"
)

func JaegerTrace(c *gin.Context) {
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)

	operationName := c.Request.Method + ":" + c.Request.URL.Path
	var span opentracing.Span
	if err == nil {
		span = opentracing.StartSpan(operationName, ext.RPCServerOption(spanCtx))
	} else {
		span = opentracing.StartSpan(operationName)
	}
	defer span.Finish()

	ext.HTTPMethod.Set(span, c.Request.Method)
	ext.HTTPUrl.Set(span, c.Request.URL.Path)
	c.Set(SpanCtxKey, opentracing.ContextWithSpan(context.Background(), span))
	c.Next()
	ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
}
