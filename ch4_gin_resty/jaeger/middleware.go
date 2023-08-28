package jaeger

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	SpanCtxKey = "SpanCtxKey"
)

func Trace(c *gin.Context) {
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	operationName := c.Request.Method + " " + c.Request.URL.Path
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
	var span opentracing.Span
	if err == nil {
		span = opentracing.StartSpan(operationName, ext.RPCServerOption(spanCtx))
	} else {
		span = opentracing.StartSpan(operationName)
		if !errors.Is(err, opentracing.ErrSpanContextNotFound) {
			log.Println("gin jaeger middleware start span err:", err)
		}
	}
	defer span.Finish()

	ext.HTTPMethod.Set(span, c.Request.Method)
	ext.HTTPUrl.Set(span, c.Request.URL.Path)
	c.Set(SpanCtxKey, opentracing.ContextWithSpan(context.Background(), span))
	c.Next()
	ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
}

func GetSpanCtx(c *gin.Context) context.Context {
	value, exists := c.Get(SpanCtxKey)
	if exists {
		return value.(context.Context)
	}
	return context.Background()
}
