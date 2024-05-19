package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/wcig/go-jaeger-example/ch4_gin_resty/jaeger"
	xhttp "github.com/wcig/go-jaeger-example/ch4_gin_resty/resty"
)

func main() {
	jaeger.Init("server_bar")
	runGinServer()
}

func runGinServer() {
	router := gin.Default()
	router.Use(jaeger.Trace)
	router.GET("/bar", barHandler)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	gracefulShutdown(srv)
}

func barHandler(c *gin.Context) {
	ctx := jaeger.GetSpanCtx(c)
	span1, spanCtx := opentracing.StartSpanFromContext(ctx, "barHandler")
	defer span1.Finish()

	time.Sleep(100 * time.Millisecond)

	resp, err := xhttp.TraceClient.R().SetContext(spanCtx).Get("http://localhost:8082/foo")
	log.Println(resp.String(), err)

	c.JSON(http.StatusOK, map[string]interface{}{"server": "bar"})
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go shutdownHook()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func shutdownHook() {
	jaeger.Close()
}
