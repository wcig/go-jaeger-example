package main

import (
	"fmt"
	"github.com/wcig/go-jaeger-example/lib/util"
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
		fmt.Println(">> server bar receive request headers:", util.ToPrettyJsonStr(r.Header))

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

	// 2025/02/23 22:15:33 debug logging disabled
	// 2025/02/23 22:15:33 Initializing logging reporter
	// 2025/02/23 22:15:33 debug logging disabled
	// >> server bar receive request headers: {
	//        "Accept-Encoding": [
	//                "gzip"
	//        ],
	//        "Uber-Trace-Id": [
	//                "733f787dd1931b3f:78381841d9c5bb2c:733f787dd1931b3f:1"
	//        ],
	//        "Uberctx-Date": [
	//                "20230831"
	//        ],
	//        "User-Agent": [
	//                "Go-http-client/1.1"
	//        ]
	// }
	// barHandler get baggage item: 20230831
	// 2025/02/23 22:15:53 Reporting span 733f787dd1931b3f:7d598b2f1e950e79:78381841d9c5bb2c:1
}
