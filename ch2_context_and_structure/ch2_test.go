package main

import (
	"context"
	"goapp/lib/tracing"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	serviceName = "ch2_context_and_structure"
)

// TestTraceWithSeparateSpans three separate span
func TestTraceWithSeparateSpans(t *testing.T) {
	tracer, closer := tracing.InitJaeger(serviceName)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	span := opentracing.StartSpan("TestTraceWithSeparateSpans")
	defer span.Finish()
	time.Sleep(100 * time.Millisecond)

	fn1 := func() {
		span1 := opentracing.StartSpan("fn1")
		defer span1.Finish()
		time.Sleep(200 * time.Millisecond)
	}

	fn2 := func() {
		span2 := opentracing.StartSpan("fn2")
		defer span2.Finish()
		time.Sleep(300 * time.Millisecond)
	}

	fn1()
	fn2()

	// Output:
	// 2023/08/27 23:13:09 debug logging disabled
	// 2023/08/27 23:13:09 Initializing logging reporter
	// 2023/08/27 23:13:09 debug logging disabled
	// 2023/08/27 23:13:10 Reporting span 356e2a54a88bed12:356e2a54a88bed12:0000000000000000:1
	// 2023/08/27 23:13:10 Reporting span 08a7383ba4652371:08a7383ba4652371:0000000000000000:1
	// 2023/08/27 23:13:10 Reporting span 3913a9e23d833dc9:3913a9e23d833dc9:0000000000000000:1
}

// TestTraceWithContext parent span, child span, follow span with context
func TestTraceWithContext(t *testing.T) {
	tracer, closer := tracing.InitJaeger(serviceName)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	span := opentracing.StartSpan("TestTraceWithContext")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	time.Sleep(100 * time.Millisecond)

	fn1 := func(ctx context.Context) {
		span1, _ := opentracing.StartSpanFromContext(ctx, "fn1")
		defer span1.Finish()
		time.Sleep(200 * time.Millisecond)
	}

	fn2 := func(ctx context.Context) {
		span2, _ := opentracing.StartSpanFromContext(ctx, "fn2")
		defer span2.Finish()
		time.Sleep(300 * time.Millisecond)
	}

	fn1(ctx)
	fn2(ctx)

	// Output:
	// 2023/08/27 23:13:36 debug logging disabled
	// 2023/08/27 23:13:36 Initializing logging reporter
	// 2023/08/27 23:13:36 debug logging disabled
	// 2023/08/27 23:13:37 Reporting span 1b9eff78cd8789b4:09b4e5d293fd8cb2:1b9eff78cd8789b4:1
	// 2023/08/27 23:13:37 Reporting span 1b9eff78cd8789b4:2e514b78457037a4:1b9eff78cd8789b4:1
	// 2023/08/27 23:13:37 Reporting span 1b9eff78cd8789b4:1b9eff78cd8789b4:0000000000000000:1
}

// TestTraceWithNestedFunc nested function
func TestTraceWithNestedFunc(t *testing.T) {
	tracer, closer := tracing.InitJaeger(serviceName)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	span := opentracing.StartSpan("TestTraceWithNestedFunc")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	time.Sleep(100 * time.Millisecond)

	fn1 := func(ctx context.Context) {
		fn11 := func(ctx context.Context) {
			span11, _ := opentracing.StartSpanFromContext(ctx, "fn11")
			defer span11.Finish()
			time.Sleep(20 * time.Millisecond)
		}

		span1, ctx1 := opentracing.StartSpanFromContext(ctx, "fn1")
		defer span1.Finish()
		time.Sleep(200 * time.Millisecond)

		fn11(ctx1)
	}

	fn2 := func(ctx context.Context) {
		span2, _ := opentracing.StartSpanFromContext(ctx, "fn2")
		defer span2.Finish()
		time.Sleep(300 * time.Millisecond)
	}

	fn1(ctx)
	fn2(ctx)

	// Output:
	// 2023/08/27 23:14:20 debug logging disabled
	// 2023/08/27 23:14:20 Initializing logging reporter
	// 2023/08/27 23:14:20 debug logging disabled
	// 2023/08/27 23:14:20 Reporting span 698c06ae9729dd74:5ccafbcc36a9cfeb:7922f2ab98a34e60:1
	// 2023/08/27 23:14:20 Reporting span 698c06ae9729dd74:7922f2ab98a34e60:698c06ae9729dd74:1
	// 2023/08/27 23:14:21 Reporting span 698c06ae9729dd74:268ece0062df9dcf:698c06ae9729dd74:1
	// 2023/08/27 23:14:21 Reporting span 698c06ae9729dd74:698c06ae9729dd74:0000000000000000:1
}
