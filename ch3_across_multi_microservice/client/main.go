package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	xhttp "github.com/wcig/go-jaeger-example/lib/http"
	"github.com/wcig/go-jaeger-example/lib/tracing"
)

// Trace a transaction across more than one microservice
func main() {
	tracer, closer := tracing.InitJaeger("client")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	span := opentracing.StartSpan("main")
	span.SetTag("main", 111)
	span.SetBaggageItem("date", "20230831")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	time.Sleep(100 * time.Millisecond)

	callBar(ctx)
	callFoo(ctx)

	// Output:
	// 2023/08/27 23:44:32 debug logging disabled
	// 2023/08/27 23:44:32 Initializing logging reporter
	// 2023/08/27 23:44:32 debug logging disabled
	// callBar resp: {"server":"bar"}
	// 2023/08/27 23:44:32 Reporting span 1d9752378fe4d05d:02543d12662022d3:1d9752378fe4d05d:1
	// callFoo resp: {"server":"foo"}
	// 2023/08/27 23:44:33 Reporting span 1d9752378fe4d05d:187fe4771a52880d:1d9752378fe4d05d:1
	// 2023/08/27 23:44:33 Reporting span 1d9752378fe4d05d:1d9752378fe4d05d:0000000000000000:1

	// 2025/02/23 22:15:52 debug logging disabled
	// 2025/02/23 22:15:52 Initializing logging reporter
	// 2025/02/23 22:15:52 debug logging disabled
	// callBar get baggage item: 20230831
	// callBar resp: {"server":"bar"}
	// 2025/02/23 22:15:53 Reporting span 733f787dd1931b3f:78381841d9c5bb2c:733f787dd1931b3f:1
	// callFoo resp: {"server":"foo"}
	// 2025/02/23 22:15:53 Reporting span 733f787dd1931b3f:2d16d78141b9ec93:733f787dd1931b3f:1
	// 2025/02/23 22:15:53 Reporting span 733f787dd1931b3f:733f787dd1931b3f:0000000000000000:1
}

func callBar(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callBar")
	span.SetTag("callBar", 222)
	defer span.Finish()

	date := span.BaggageItem("date")
	fmt.Println("callBar get baggage item:", date)

	url := "http://localhost:8081/bar"
	method := http.MethodGet
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, method)
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		panic(err)
	}

	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("callBar resp:", string(resp))
}

func callFoo(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callFoo")
	span.SetTag("callFoo", 333)
	defer span.Finish()

	url := "http://localhost:8082/foo"
	method := http.MethodGet
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, method)
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		panic(err)
	}

	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("callFoo resp:", string(resp))
}
