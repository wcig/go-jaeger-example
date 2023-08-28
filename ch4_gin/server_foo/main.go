package main

import (
	"context"
	"goapp/ch4_gin/jaeger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	jaeger.Init("server_foo")
	runGinServer()
}

func runGinServer() {
	router := gin.Default()
	router.Use(jaeger.Trace)
	router.GET("/foo", fooHandler)
	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	gracefulShutdown(srv)
}

func fooHandler(c *gin.Context) {
	ctx := jaeger.GetSpanCtx(c)
	span, _ := opentracing.StartSpanFromContext(ctx, "fooHandler")
	defer span.Finish()

	time.Sleep(200 * time.Millisecond)
	c.JSON(http.StatusOK, map[string]interface{}{"server": "foo"})
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
