package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	v1 "github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/api/service/v1"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/conf"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/server"
)

type Config struct {
	GRPCPort string
	DataStoreDBHost string
	DataStoreDBUser string
	DataStoreDBPassword string
	DataStoreDBSchema string
}
var cfg Config
func init()  {
	cfg=Config{}
	flag.StringVar(&cfg.GRPCPort,"port",conf.Port,"gRPC port to bind")
	flag.StringVar(&cfg.DataStoreDBHost,"db-host",conf.DbHost,"db host")
	flag.StringVar(&cfg.DataStoreDBUser,"db-user",conf.DbUser,"db-user")
	flag.StringVar(&cfg.DataStoreDBPassword,"db-password",conf.DbPassword,"db-password")
	flag.StringVar(&cfg.DataStoreDBSchema,"db-schema",conf.DbSchema,"db-schema")
	fmt.Println("init"+cfg.GRPCPort)
	flag.Parse()
}

func RunServer() error {
	ctx :=context.Background()

	if len(cfg.GRPCPort)==0{
		return fmt.Errorf("invalid TCP port for gRPC server：%s",cfg.GRPCPort)
	}
	param:="parseTime=true"
	dsn :=fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.DataStoreDBUser, cfg.DataStoreDBPassword, cfg.DataStoreDBHost, cfg.DataStoreDBSchema, param)
	db,err:=sql.Open("mysql",dsn)
	if err!=nil{
		return fmt.Errorf("连接数据失败:%v",err)
	}
	defer db.Close()

	v1API:= v1.NewToDoServiceServer(db)
	return server.RunServer(ctx,v1API,cfg.GRPCPort)
}