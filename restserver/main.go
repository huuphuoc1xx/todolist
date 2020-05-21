package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"todolist/controller"

	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
	var Auth = rand.Int()
	if Auth%2 == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauth"})
		c.Abort()
	}
	fmt.Printf("%s", Auth)
	c.Next()

	return
}
func main() {
	controller.InitGRPC()
	g := gin.Default()
	g.Use(Authorize)
	g.GET("/list-todo", controller.GetListToDo)
	g.GET("/todo-by-id/:id", controller.GetToDo)

	g.POST("/create-task", controller.CreateToDo)
	g.PUT("/update-task", controller.UpdateToDo)
	g.DELETE("/delete-task/:id", controller.DeleteToDo)
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
