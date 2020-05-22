package cmd

import (
	"context"
	"fmt"
	"todolist/TodolistServer/grpc"
	"todolist/TodolistServer/server"
	"todolist/config"
)

func RunServer() error {
	ctx := context.Background()
	var Dbinfo config.DBInfo
	Dbinfo.SetTodoDB()
	db, err := Dbinfo.GetDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	API := server.NewToDoServer(db)
	return grpc.RunServer(ctx, API, "2000")
}
