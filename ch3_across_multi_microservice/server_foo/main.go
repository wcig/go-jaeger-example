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
	tracer, closer := tracing.InitJaeger("server_foo")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	http.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(">> server foo receive request headers:", util.ToPrettyJsonStr(r.Header))

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

	// 2025/02/23 22:15:27 debug logging disabled
	// 2025/02/23 22:15:27 Initializing logging reporter
	// 2025/02/23 22:15:27 debug logging disabled
	// >> server foo receive request headers: {
	//        "Accept-Encoding": [
	//                "gzip"
	//        ],
	//        "Uber-Trace-Id": [
	//                "733f787dd1931b3f:2d16d78141b9ec93:733f787dd1931b3f:1"
	//        ],
	//        "Uberctx-Date": [
	//                "20230831"
	//        ],
	//        "User-Agent": [
	//                "Go-http-client/1.1"
	//        ]
	// }
	// fooHandler get baggage item: 20230831
	// 2025/02/23 22:15:53 Reporting span 733f787dd1931b3f:6fec2a3099e24af3:2d16d78141b9ec93:1
}
