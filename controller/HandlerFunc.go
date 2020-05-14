package controller

import (
	"net/http"
	"strconv"
	"todolist/config"

	"github.com/gin-gonic/gin"
)

func GetListToDo(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pagenumber"})
		return
	}

	tag := ctx.Query("tag")
	if len(tag) == 0 {
		tag = "*"
	}
	todos := getTaskByPage(tag, page)
	ctx.JSON(200, todos)
}

func CreateToDo(ctx *gin.Context) {

	var todo config.ToDo
	ctx.BindJSON(&todo)
	id := createTask(todo)
	if id != 0 {
		ctx.JSON(200, gin.H{"ID": id})
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

	todo := getTaskById(id)
	if todo == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ctx.JSON(200, todo)
}

func UpdateToDo(ctx *gin.Context) {

	var todo config.ToDo
	ctx.BindJSON(&todo)
	id := updateTask(todo)
	if id != 0 {
		ctx.JSON(200, gin.H{"ID": id})
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
	id = deleteTask(id)
	if id != 0 {
		ctx.JSON(200, gin.H{"ID": id})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something was wrong!!"})
	}
}
