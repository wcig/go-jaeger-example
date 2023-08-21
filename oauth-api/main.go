package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var (
	rc = resty.New()
)

func main() {
	InitJaeger()
	InitGin()
}

func InitGin() {
	e := gin.Default()
	e.Use(JaegerTrace)
	e.POST("/oauth2/client/delete", deleteClientHandler)
	err := e.Run(":28080")
	if err != nil {
		panic(err)
	}
}

func deleteClientHandler(c *gin.Context) {
	spanCtxVal, _ := c.Get(SpanKey)
	spanCtx := spanCtxVal.(context.Context)

	checkPermission(spanCtx)
	deleteOAuthClient(spanCtx)
	deleteOAuthCode(spanCtx)
	deleteOAuthToken(spanCtx)
	reportAuditLog(spanCtx)

	c.JSON(http.StatusOK, map[string]interface{}{"code": 0})
}

func checkPermission(ctx context.Context) {
	userID := "0a15f6"
	span, _ := opentracing.StartSpanFromContext(ctx, "check_permission")
	span.SetTag("user_id", userID)
	defer span.Finish()

	url := "http://localhost:28081/permission/check"
	req := rc.R()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		fmt.Println(err)
	}

	resp, err := req.SetBody(map[string]interface{}{"user_id": userID}).Post(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = resp
}

func deleteOAuthClient(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "delete_oauth_client")
	span.SetTag("client_id", "111222333")
	defer span.Finish()

	time.Sleep(20 * time.Millisecond)
}

func deleteOAuthCode(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "delete_oauth_code")
	span.SetTag("code", "111111")
	defer span.Finish()

	time.Sleep(30 * time.Millisecond)
}

func deleteOAuthToken(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "delete_oauth_token")
	span.SetTag("token", "a1b2c3")
	defer span.Finish()

	time.Sleep(40 * time.Millisecond)
}

func reportAuditLog(ctx context.Context) {
	// todo
}
