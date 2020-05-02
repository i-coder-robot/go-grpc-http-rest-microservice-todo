package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const (
	apiVersion="v1"
)

type ToDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) *ToDoServiceServer{
	return &ToDoServiceServer{db: db}
}

func (s *ToDoServiceServer) checkAPI(api string) error {
	if len(api)>0{
		if apiVersion!=api{
			msg :="unsupported API version:service implements API version "+apiVersion+",but given "+api
			return status.Error(codes.Unimplemented,msg)
		}
	}
	return nil
}

func (s *ToDoServiceServer) connect(ctx context.Context) (*sql.Conn,error){
	c,err:=s.db.Conn(ctx)
	if err!=nil{
		return nil,status.Error(codes.Unknown,"连接数据库失败."+err.Error())
	}
	return c,nil
}

func (s *ToDoServiceServer)  Create(ctx context.Context,req *v1.CreateRequest)(*v1.CreateResponse,error){
	if err:=s.checkAPI(req.Api);err!=nil{
		return nil, err
	}
	c,err:=s.connect(ctx)
	if err!=nil{
		return nil,err
	}
	defer c.Close()
	reminder,err := ptypes.Timestamp(req.ToDo.Reminder)
	if err!=nil{
		return nil,status.Error(codes.InvalidArgument,"参数错误:"+err.Error())
	}
	res,err:=c.ExecContext(ctx,"INSERT INTO ToDo(`Title`,`Description`,`Reminder`) VALUES(?,?,?)",req.ToDo.Title,req.ToDo.Description,reminder)
	if err!=nil{
		return nil,status.Error(codes.Unknown,"添加 ToDo失败"+err.Error())
	}
	id,err:=res.LastInsertId()
	if err!=nil{
		return nil,status.Error(codes.Unknown,"获取最新ID失败"+err.Error())
	}
	return &v1.CreateResponse{Api: apiVersion,Id: id},nil
}

func (s *ToDoServiceServer) Read(ctx context.Context,req *v1.ReadRequest) (*v1.ReadResponse,error){
	if err:=s.checkAPI(req.Api);err!=nil{
		return nil,err
	}
	c,err:=s.connect(ctx)
	if err!=nil{
		return nil,err
	}
	defer c.Close()

	rows,err:=c.QueryContext(ctx, "SELECT `ID`,`Title`,`Description`,`Reminder` FROM ToDo WHERE `ID`=? ",req.Id)
	if err!=nil{
		return nil,status.Error(codes.Unknown,"查询失败:"+err.Error())
	}
	defer rows.Close()
	if !rows.Next(){
		if err:=rows.Err();err!=nil{
			return nil,status.Error(codes.Unknown,"失败获取数据:"+err.Error())
		}
		return nil,status.Error(codes.NotFound,fmt.Sprintf("ID='%d'找不到",req.Id))
	}
	var td v1.ToDo
	var reminder time.Time
	if err:=rows.Scan(&td.Id,&td.Title,&td.Description,&reminder);err!=nil{
		return nil,status.Error(codes.Unknown,"查找数据失败:"+err.Error())
	}
	if rows.Next(){
		return nil,status.Error(codes.Unknown,fmt.Sprintf("查找到多条数据ID='%d'",req.Id))
	}
	return &v1.ReadResponse{Api: apiVersion,ToDo: &td},nil
}

func (s *ToDoServiceServer) Update(ctx context.Context,req *v1.UpdateRequest) (*v1.UpdateResponse,error){
	if err:=s.checkAPI(req.Api);err!=nil{
		return nil,err
	}
	c,err:=s.connect(ctx)
	if err != nil {
		return nil,err
	}
	defer c.Close()

	reminder,err:=ptypes.Timestamp(req.ToDo.Reminder)
	if err!=nil{
		return nil,status.Error(codes.InvalidArgument,"reminder参数无效")
	}
	res,err:=c.ExecContext(ctx,"UPDATE ToDo SET `Title`=?,`Description`=?,`Reminder`=? WHERE `ID`=?",
		req.ToDo.Title,req.ToDo.Description,reminder,req.ToDo.Id)
	if err != nil {
		return nil,status.Error(codes.Unknown,"更新失败:"+err.Error())
	}
	rows,err:=res.RowsAffected()
	if err!=nil{
		return nil,status.Error(codes.Unknown,"失败有效的行更新"+err.Error())
	}
	if rows==0{
		msg := "ID="+strconv.FormatInt(req.ToDo.Id,10)+"找不到"
		return nil,status.Error(codes.NotFound,msg)
	}
	return &v1.UpdateResponse{
		Api: apiVersion,
		Updated: rows,
	},nil
}

func (s *ToDoServiceServer) Delete(ctx context.Context,req *v1.DeleteRequest) (*v1.DeleteResponse,error){
	if err:=s.checkAPI(req.Api);err!=nil{
		return nil,err
	}
	c,err:=s.connect(ctx)
	if err != nil {
		return nil,err
	}
	defer c.Close()
	res,err:=c.ExecContext(ctx,"DELETE FROM ToDo where 'ID'=?",req.Id)
	if err!=nil{
		return nil,status.Error(codes.Unknown,"删除失败:"+err.Error())
	}
	rows,err:=res.RowsAffected()
	if err!=nil{
		return nil,status.Error(codes.Unknown,"失败更新行失败:"+err.Error())
	}
	if rows==0{
		return nil,status.Error(codes.NotFound,fmt.Sprintf("ID='%d'未找到",req.Id))
	}
	return &v1.DeleteResponse{
		Api: req.Api,
		Deleted: rows,
	},nil
}

func (s *ToDoServiceServer) ReadAll(ctx context.Context,req *v1.ReadAllRequest) (*v1.ReadAllResponse,error){
	if err:= s.checkAPI(req.Api);err!=nil{
		return nil,err
	}
	c,err:=s.connect(ctx)
	if err!=nil{
		return nil,err
	}
	defer c.Close()

	rows,err := c.QueryContext(ctx,"SELECT `ID`,`Title`,`Description`,`Reminder` FROM ToDo")
	if err!=nil{
		return nil,status.Error(codes.Unknown,"查询失败:"+err.Error())
	}
	defer rows.Close()
	var reminder time.Time
	list := []*v1.ToDo{}
	for rows.Next(){
		td:=new(v1.ToDo)
		if err:=rows.Scan(&td.Id,&td.Title,&td.Description,&reminder);err!=nil{
			return nil,status.Error(codes.Unknown,"查询失败"+err.Error())
		}
		td.Reminder,err=ptypes.TimestampProto(reminder)
		if err!=nil{
			return nil,status.Error(codes.Unknown,"reminder 无效:"+err.Error())
		}
		list=append(list,td)
	}
	if err:=rows.Err();err!=nil{
		return nil,status.Error(codes.Unknown,"获取数据失败:"+err.Error())
	}
	return &v1.ReadAllResponse{
		Api: apiVersion,
		ToDos: list,
	},nil
}