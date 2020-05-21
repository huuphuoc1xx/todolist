package server

import (
	"context"
	"database/sql"
	"time"
	todoserver "todolist/TodolistServer/proto"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/go-sql-driver/mysql"
)

const pagesize = 2

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServer(db *sql.DB) todoserver.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	return s.db.Conn(ctx)
}

func (s *toDoServiceServer) GetByTag(ctx context.Context, req *todoserver.GetByTagRequest) (*todoserver.GetByTagResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := c.QueryContext(ctx, "Select * from todolist where Tag like ? limit ? offset ?", req.GetTag(), pagesize, (req.GetPage()-1)*pagesize)
	if err != nil {
		return nil, err
	}
	list := make([]*todoserver.ToDo, 0)
	for rows.Next() {
		todo := new(todoserver.ToDo)
		var CreateTime string
		if err := rows.Scan(&todo.Username, &todo.Id, &todo.Title, &todo.Tag, &todo.Description, &CreateTime); err != nil {
			return nil, err
		}

		formatTime := "2006-01-02 15:04:05"
		createTime, _ := time.Parse(formatTime, CreateTime)
		todo.CreateTime, err = ptypes.TimestampProto(createTime)
		if err != nil {
			return nil, status.Error(codes.Unknown, "CreateTime field has invalid format-> "+err.Error())
		}
		list = append(list, todo)
	}

	return &todoserver.GetByTagResponse{
		Count: int64(len(list)),
		Todos: list,
	}, nil
}

func (s *toDoServiceServer) Create(ctx context.Context, req *todoserver.ToDoRequest) (*todoserver.IdResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	todo := req.GetTodo()
	res, err := c.ExecContext(ctx, "Insert into todolist(Username,Title,Tag,Description) value (?,?,?,?)",
		todo.GetUsername(), todo.GetTitle(), todo.GetTag(), todo.GetDescription())

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &todoserver.IdResponse{Id: id}, nil
}

func (s *toDoServiceServer) Update(ctx context.Context, req *todoserver.ToDoRequest) (*todoserver.IdResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	todo := req.GetTodo()
	_, err = c.ExecContext(ctx, "Update todolist Set Title=?,Tag=?,Description=? where ID=?",
		todo.GetTitle(), todo.GetTag(), todo.GetDescription(), todo.GetId())
	if err != nil {
		return nil, err
	}
	return &todoserver.IdResponse{Id: todo.GetId()}, nil
}

func (s *toDoServiceServer) GetById(ctx context.Context, req *todoserver.GetByIdRequest) (*todoserver.GetByIdResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	row := c.QueryRowContext(ctx, "Select * from todolist where ID=?", req.GetId())

	var todo todoserver.ToDo
	var CreateTime string
	if err := row.Scan(&todo.Username, &todo.Id, &todo.Title, &todo.Tag, &todo.Description, &CreateTime); err != nil {
		return nil, err
	}

	formatTime := "2006-01-02 15:04:05"
	createTime, _ := time.Parse(formatTime, CreateTime)
	todo.CreateTime, err = ptypes.TimestampProto(createTime)
	if err != nil {
		return nil, status.Error(codes.Unknown, "CreateTime field has invalid format-> "+err.Error())
	}
	return &todoserver.GetByIdResponse{Todo: &todo}, nil
}

func (s *toDoServiceServer) Delete(ctx context.Context, req *todoserver.DeleteRequest) (*todoserver.DeleteResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	res, err := c.ExecContext(ctx, "Delete from todolist where ID=?", req.GetId)
	if err != nil {
		return nil, err
	}

	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		return &todoserver.DeleteResponse{Success: 0}, nil
	}
	return &todoserver.DeleteResponse{Success: 1}, nil
}
