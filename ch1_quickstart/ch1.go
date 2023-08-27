package main

import (
	"goapp/lib/tracing"
	"time"

	"github.com/opentracing/opentracing-go/log"
)

// 快速开始
func main() {
	tracer, closer := tracing.InitJaeger("ch1_quickstart")
	defer closer.Close()

	// span
	span := tracer.StartSpan("part1")
	defer span.Finish()
	time.Sleep(123 * time.Millisecond)

	// tag
	span.SetTag("my_key", "my_value")

	// log
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", "hello jaeger!"),
	)
	span.LogKV("event", "println")

	// Output:
	// 2023/08/27 23:10:01 debug logging disabled
	// 2023/08/27 23:10:01 Initializing logging reporter
	// 2023/08/27 23:10:01 debug logging disabled
	// 2023/08/27 23:10:02 Reporting span 03e2f3eb3f4d9303:03e2f3eb3f4d9303:0000000000000000:1
}
