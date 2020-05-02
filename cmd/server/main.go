package main

import (
	"fmt"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/cmd"
	"os"
)



func main() {
	if err:=cmd.RunServer();err!=nil{
		fmt.Fprintf(os.Stderr,"%v\n",err)
		os.Exit(1)
	}
}
