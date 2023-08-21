package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	InitJaeger()
	InitGin()
}

func InitGin() {
	e := gin.Default()
	e.Use(JaegerTrace)
	e.POST("/permission/check", checkPermissionHandler)
	err := e.Run(":28081")
	if err != nil {
		panic(err)
	}
}

func checkPermissionHandler(c *gin.Context) {
	spanCtxVal, _ := c.Get(SpanKey)
	spanCtx := spanCtxVal.(context.Context)

	span, _ := opentracing.StartSpanFromContext(spanCtx, "check_permission")
	defer span.Finish()

	time.Sleep(100 * time.Millisecond)
}
