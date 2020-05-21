package server

import (
	"context"
	"database/sql"
	"time"
	profile "todolist/ProfileServer/proto"

	"github.com/golang/profilebuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/go-sql-driver/mysql"
)

const pagesize = 2

type profileServiceServer struct {
	db *sql.DB
}

func NewProfileServer(db *sql.DB) profile.ProfileServiceServer {
	return &profileServiceServer{db: db}
}

func (s *profileServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	return s.db.Conn(ctx)
}

func (s *profileServiceServer) Register(ctx context.Context, req *profile.RegisterRequest) (*profile.RegisterResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	info := req.GetProfile()
	res, err := c.ExecContext(ctx, "Insert into profile(Username,Password,Name,Email,Phone) value (?,?,?,?)",
		info.GetUsername(), info.GetPassword(), info.GetName(), info.GetEmail(), info.GetPhone())

	if err != nil {
		return &profile.RegisterResponse{Success: false}, err
	}
	return &profile.IdResponse{Success: true}, nil
}

func (s *profileServiceServer) Update(ctx context.Context, req *profile.UpdateRequest) (*profile.UpdateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	info := req.GetProfile()
	res, err := c.ExecContext(ctx, "Update profile Set Password=?, Name=?, Email=?,Phone=? where Username=?",
		req.GetPassword(), info.GetName(), info.GetEmail(), info.GetPhone(), req.GetUsername())

	r := res.RowsAffected()
	if err != nil {
		return &profile.RegisterResponse{Success: false}, err
	}
	return &profile.IdResponse{Success: true}, nil
}

func (s *profileServiceServer) GetProfile(ctx context.Context, req *profile.GetProfileRequest) (*profile.Profile, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	row := c.QueryRowContext(ctx, "Select Name,Email,Phone from Profile where Username=?", req.GetId())

	var todo profile.Profile
	var CreateTime string
	if err := row.Scan(&todo.Username, &todo.Id, &todo.Title, &todo.Tag, &todo.Description, &CreateTime); err != nil {
		return nil, err
	}

	formatTime := "2006-01-02 15:04:05"
	createTime, _ := time.Parse(formatTime, CreateTime)
	todo.CreateTime, err = ptypes.Timestampprofile(createTime)
	if err != nil {
		return nil, status.Error(codes.Unknown, "CreateTime field has invalid format-> "+err.Error())
	}
	return &profile.GetByIdResponse{Todo: &todo}, nil
}

func (s *profileServiceServer) Delete(ctx context.Context, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	res, err := c.ExecContext(ctx, "Delete from todolist where ID=?", req.GetId)
	if err != nil {
		return nil, err
	}

	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		return &profile.DeleteResponse{Success: 0}, nil
	}
	return &profile.DeleteResponse{Success: 1}, nil
}
