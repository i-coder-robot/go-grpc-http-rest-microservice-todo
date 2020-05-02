package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/api/proto/v1"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/conf"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "127.0.0.1:"+conf.Port, "gRPC server in format host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("服务端,连不上啊:%v", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatal("创建失败%v", err)
	}
	log.Printf("Create result%v", res1)
	id := res1.Id
	fmt.Sprintf("%v",res1)

	req2 := v1.ReadRequest{Api: apiVersion, Id: id}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatal("Read failed %v", err)
	}
	log.Printf("Read reslut %v", res2)

	req3 :=v1.UpdateRequest{
		Api:                  apiVersion,
		ToDo:                 &v1.ToDo{
			Id:                   res2.ToDo.Id,
			Title:                res2.ToDo.Title,
			Description:          res2.ToDo.Description+" updated",
			Reminder:             res2.ToDo.Reminder,
		},
	}
	res3,err:=c.Update(ctx,&req3)
	if err!=nil{
		log.Fatal("更新失败%v",err)
	}
	log.Printf("Update result %v",res3)

	req4 := v1.ReadAllRequest{
		Api: apiVersion,
	}

	res4,err:=c.ReadAll(ctx,&req4)
	if err!=nil{
		log.Fatal("ReadAll 失败%v",err)
	}
	log.Printf("ReadAll result %v",res4)

	req5 :=v1.DeleteRequest{
		Api: apiVersion,
		Id: id,
	}
	res5,err:=c.Delete(ctx,&req5)
	if err!=nil{
		log.Fatal("删除失败%v",err)
	}
	log.Printf("删除的结果%v",res5)
}
