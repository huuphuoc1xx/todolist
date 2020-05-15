package main

import (
	"log"
	"todolist/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/list-todo", controller.GetListToDo)
	g.GET("/todo-by-id/:id", controller.GetToDo)

	g.POST("/create-task", controller.CreateToDo)
	g.PUT("/update-task", controller.UpdateToDo)
	g.DELETE("/delete-task/:id", controller.DeleteToDo)
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
