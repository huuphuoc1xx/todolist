package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"todolist/config"

	"github.com/gin-contrib/sessions"

	profileproto "todolist/ProfileServer/proto"
	todoproto "todolist/TodolistServer/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var toDoClient todoproto.ToDoServiceClient
var profileClient profileproto.ProfileServiceClient

func InitGRPC() {
	conn, err := grpc.Dial("localhost:2000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	toDoClient = todoproto.NewToDoServiceClient(conn)

	conn, err = grpc.Dial("localhost:2001", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}
	profileClient = profileproto.NewProfileServiceClient(conn)
}

func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("Auth") == "true" {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Already login!!!"})
	}

	var loginreq profileproto.LoginRequest
	ctx.BindJSON(&loginreq)

	_, err := profileClient.Login(ctx, &loginreq)
	if err != nil {
		ctx.JSON(200, gin.H{"error": err.Error()})
		return
	}

	session.Set("Auth", "true")
	session.Set("AuthUser", loginreq.GetUsername())
	session.Save()

	ctx.JSON(200, "login successful")
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set("Auth", false)
	session.Save()

	ctx.JSON(200, "logout successful")
}

func Register(ctx *gin.Context) {
	var req profileproto.RegisterRequest
	ctx.BindJSON(&req)

	_, err := profileClient.Register(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, "Register Successful")
}

func UpdateProfile(ctx *gin.Context) {
	var req profileproto.UpdateRequest
	ctx.BindJSON(&req)

	_, err := profileClient.Update(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, "Update Successful")
}

func GetProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	var req profileproto.GetProfileRequest
	ctx.BindJSON(&req)

	req.Username = uname

	res, err := profileClient.GetProfile(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, res)
}
func GetListToDo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pagenumber"})
		return
	}

	tag := ctx.Query("tag")
	if len(tag) == 0 {
		tag = "%"
	}
	todos, err := toDoClient.GetByTag(ctx,
		&todoproto.GetByTagRequest{
			Username: uname,
			Page:     int64(page),
			Tag:      tag,
		})
	if err != nil {
		ctx.JSON(200, err)
		return
	}
	ctx.JSON(200, todos.GetTodos())
}

func CreateToDo(ctx *gin.Context) {

	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	var todo config.ToDo
	ctx.BindJSON(&todo)

	id, err := toDoClient.Create(ctx, &todoproto.ToDoRequest{
		Todo: &todoproto.ToDo{
			Title:       todo.Title,
			Tag:         todo.Tag,
			Username:    uname,
			Description: todo.Description,
		},
	})
	if err == nil {
		ctx.JSON(200, gin.H{"ID": id.GetId()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
}

func GetToDo(ctx *gin.Context) {

	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	todo, err := toDoClient.GetById(ctx, &todoproto.GetByIdRequest{
		Username: uname,
		Id:       int64(id),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
		return
	}
	if todo == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ctx.JSON(200, todo.GetTodo())
}

func UpdateToDo(ctx *gin.Context) {

	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	var todo config.ToDo
	ctx.BindJSON(&todo)

	result, err := toDoClient.Update(ctx, &todoproto.ToDoRequest{
		Todo: &todoproto.ToDo{
			Id:          todo.Id,
			Title:       todo.Title,
			Tag:         todo.Tag,
			Username:    uname,
			Description: todo.Description,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
	if result != nil {
		ctx.JSON(200, gin.H{"ID": result.GetId()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
}

func DeleteToDo(ctx *gin.Context) {

	session := sessions.Default(ctx)
	uname := fmt.Sprint(session.Get("AuthUser"))

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Invalid Id"})
		return
	}
	result, err := toDoClient.Delete(ctx, &todoproto.DeleteRequest{
		Username: uname,
		Id:       int64(id),
	})
	if err == nil && result.GetSuccess() == 1 {
		ctx.JSON(200, gin.H{"ID": id})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
}
