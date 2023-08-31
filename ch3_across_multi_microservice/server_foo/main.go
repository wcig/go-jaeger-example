package main

import (
	"fmt"
	"goapp/lib/tracing"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	tracer, closer := tracing.InitJaeger("server_foo")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	http.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			panic(err)
		}

		span := opentracing.StartSpan("foo", ext.RPCServerOption(spanCtx))
		span.SetTag("foo", 222222)
		defer span.Finish()

		date := span.BaggageItem("date")
		fmt.Println("fooHandler get baggage item:", date)

		time.Sleep(300 * time.Millisecond)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(`{"server":"foo"}`))
	})
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		panic(err)
	}

	// Output:
	// 2023/08/27 23:44:29 debug logging disabled
	// 2023/08/27 23:44:29 Initializing logging reporter
	// 2023/08/27 23:44:29 debug logging disabled
	// 2023/08/27 23:44:33 Reporting span 1d9752378fe4d05d:2ae8b39babac786b:187fe4771a52880d:1
}
