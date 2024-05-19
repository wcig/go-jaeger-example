package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/wcig/go-jaeger-example/lib/tracing"
)

func main() {
	tracer, closer := tracing.InitJaeger("server_bar")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	http.HandleFunc("/bar", func(rw http.ResponseWriter, r *http.Request) {
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			panic(err)
		}

		span := opentracing.StartSpan("bar", ext.RPCServerOption(spanCtx))
		span.SetTag("bar", 111111)
		defer span.Finish()

		date := span.BaggageItem("date")
		fmt.Println("barHandler get baggage item:", date)

		time.Sleep(200 * time.Millisecond)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(`{"server":"bar"}`))
	})
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}

	// Output:
	// 2023/08/27 23:44:27 debug logging disabled
	// 2023/08/27 23:44:27 Initializing logging reporter
	// 2023/08/27 23:44:27 debug logging disabled
	// 2023/08/27 23:44:32 Reporting span 1d9752378fe4d05d:1c6b3a92eca8aae0:02543d12662022d3:1
}
