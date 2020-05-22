package profileserver

import (
	"context"
	"database/sql"
	"fmt"
	profile "todolist/ProfileServer/proto"

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
	_, err = c.ExecContext(ctx, "Insert into Profile(Username,Password,Name,Email,Phone) value (?,?,?,?,?)",
		req.GetUsername(), req.GetPassword(), info.GetName(), info.GetEmail(), info.GetPhone())

	if err != nil {
		return &profile.RegisterResponse{Success: false}, err
	}
	return &profile.RegisterResponse{Success: true}, nil
}

func (s *profileServiceServer) Update(ctx context.Context, req *profile.UpdateRequest) (*profile.UpdateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	info := req.GetProfile()
	res, err := c.ExecContext(ctx, "Update Profile Set Password=?, Name=?, Email=?,Phone=? where Username=?",
		req.GetPassword(), info.GetName(), info.GetEmail(), info.GetPhone(), req.GetUsername())

	r, err := res.RowsAffected()
	if r == 1 && err != nil {
		return &profile.UpdateResponse{Success: false}, err
	}
	return &profile.UpdateResponse{Success: true}, nil
}

func (s *profileServiceServer) GetProfile(ctx context.Context, req *profile.GetProfileRequest) (*profile.Profile, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	row := c.QueryRowContext(ctx, "Select Name,Email,Phone from Profile where Username=?", req.GetUsername())

	var result profile.Profile
	if err := row.Scan(&result.Name, &result.Email, &result.Phone); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *profileServiceServer) Login(ctx context.Context, req *profile.LoginRequest) (*profile.LoginResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	row := c.QueryRowContext(ctx, "Select Password from Profile where Username=?", req.GetUsername())

	var Pass string
	err = row.Scan(&Pass)

	if err != nil {
		return &profile.LoginResponse{Success: false}, err
	}

	if Pass != req.GetPassword() {
		return &profile.LoginResponse{Success: false}, fmt.Errorf("Invalid Password")
	}
	return &profile.LoginResponse{Success: true}, nil
}
