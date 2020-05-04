package rest

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/api/proto/v1"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "127.0.0.1:"+grpcPort, opts); err != nil {
		log.Fatalf("启动 HTTP 网关错误: %v", err)
	}
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// 使用了 Ctrl + c 就会处理
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
	log.Println("启动 HTTP/REST 网关...")
	return srv.ListenAndServe()
}
