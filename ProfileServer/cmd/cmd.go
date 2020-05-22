package cmd

import (
	"context"
	"fmt"
	"todolist/ProfileServer/grpc"
	profileserver "todolist/ProfileServer/server"
	"todolist/config"
)

func RunServer() error {
	ctx := context.Background()
	var Dbinfo config.DBInfo
	Dbinfo.SetProfileDB()
	db, err := Dbinfo.GetDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	API := profileserver.NewProfileServer(db)
	return grpc.RunServer(ctx, API, "2001")
}
