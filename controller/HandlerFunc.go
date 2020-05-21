package controller

import (
	"log"
	"net/http"
	"strconv"
	"todolist/config"

	proto "todolist/TodolistServer/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var client proto.ToDoServiceClient

func InitGRPC() {
	conn, err := grpc.Dial("localhost:2000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewToDoServiceClient(conn)
}

func GetListToDo(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pagenumber"})
		return
	}

	tag := ctx.Query("tag")
	if len(tag) == 0 {
		tag = "%"
	}
	todos, err := client.GetByTag(ctx,
		&proto.GetByTagRequest{
			Page: int64(page),
			Tag:  tag,
		})
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx.JSON(200, todos.GetTodos())
}

func CreateToDo(ctx *gin.Context) {

	var todo config.ToDo
	ctx.BindJSON(&todo)
	id, err := client.Create(ctx, &proto.ToDoRequest{
		Todo: &proto.ToDo{
			Title:       todo.Title,
			Tag:         todo.Tag,
			Username:    todo.Username,
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	todo, err := client.GetById(ctx, &proto.GetByIdRequest{
		Id: int64(id),
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

	var todo config.ToDo
	ctx.BindJSON(&todo)
	result, err := client.Update(ctx, &proto.ToDoRequest{
		Todo: &proto.ToDo{
			Title:       todo.Title,
			Tag:         todo.Tag,
			Username:    todo.Username,
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Invalid Id"})
		return
	}
	result, err := client.Delete(ctx, &proto.DeleteRequest{
		Id: int64(id),
	})
	if err == nil && result.GetSuccess() == 1 {
		ctx.JSON(200, gin.H{"ID": id})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
}
