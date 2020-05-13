package main

import (
	"context"
	"database/sql"
	"log"
	"time"
	proto "todolist/pkg/api/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	db *sql.DB
}

const apiVersion = "v1"

func main() {
}

func (s *server) connect(ctx context.Context) (sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, status.Error(codes.Unknown, "Err:"+err.Error())
	}

	return c, nil
}

func (s *server) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "Phien ban API %s khong duoc ho tro, Ban can dung phien ban API %s", api, apiVersion)
		}
	}
	return nil
}
func (s *server) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	if err != s.checkAPI(request.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	todo := request.GetTodo()

	result, err := s.db.ExecContext(ctx, "insert into todolist(Title,Description,InsertTime) values(?,?,?)", todo.GetTitle(), todo.GetDescription(), time.Now())
	if err != nil {
		return nil, err
	}

	return result.LastInsertId(), nil
}

func (s *server) Update(ctx context.Context, request *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	return nil, nil
}
func (s *server) Delete(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	return nil, nil
}
func (s *server) Read(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
	return nil, nil
}
func (s *server) ReadAll(ctx context.Context, request *proto.AllRequest) (*proto.AllResponse, error) {
	return nil, nil
}
