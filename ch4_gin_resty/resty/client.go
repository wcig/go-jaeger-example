package resty

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	spanKey = "span"
)

var (
	TraceClient = NewTraceClient()
)

func NewTraceClient() *resty.Client {
	client := resty.New()
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		ctx := r.Context()
		operationName := r.Method + " " + r.URL
		span, _ := opentracing.StartSpanFromContext(ctx, operationName)
		r.SetContext(context.WithValue(ctx, spanKey, span))
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, r.URL)
		ext.HTTPMethod.Set(span, r.Method)
		err := span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)
		return err
	})
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		ctx := r.Request.Context()
		if value := ctx.Value(spanKey); value != nil {
			span := value.(opentracing.Span)
			ext.HTTPStatusCode.Set(span, uint16(r.StatusCode()))
			span.Finish()
		}
		return nil
	})
	return client
}
