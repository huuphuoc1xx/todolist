package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"todolist/grpcserver/grpc"
	server "todolist/grpcserver/server"
)

const (
	Username = "gsMg5DbNgQ"
	Password = "32F8Gb0lfr"
	url      = "remotemysql.com"
	port     = "3306"
	database = "gsMg5DbNgQ"

	pagesize = 2
)

func RunServer() error {
	ctx := context.Background()

	db, err := sql.Open("mysql", Username+":"+Password+"@tcp"+"("+url+":"+port+")/"+database)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	API := server.NewToDoServer(db)
	return grpc.RunServer(ctx, API, "2000")
}
