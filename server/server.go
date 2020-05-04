package server

import (
	"context"
	v1 "github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/api/proto/v1"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/cmd/middleware"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/logger"
	"google.golang.org/grpc"
	"net"
	"os"
)

func RunServer(ctx context.Context,v1API v1.ToDoServiceServer,port string) error{
	listen,err:=net.Listen("tcp",":"+port)
	if err!=nil{
		return err
	}

	opts:=[]grpc.ServerOption{}
	opts = middleware.AddLogging(logger.Log,opts)

	server :=grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server,v1API)
	c:=make(chan os.Signal,1)
	go func() {
		for range c{
			logger.Log.Warn("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
