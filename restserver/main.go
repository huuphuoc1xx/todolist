package main

import (
	"log"
	"net/http"
	"todolist/restserver/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("Auth") != "true" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauth"})
		ctx.Abort()
		return
	}
	ctx.Next()
}
func main() {
	controller.InitGRPC()
	g := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	g.POST("/login", controller.Login)
	g.POST("/register", controller.Register)

	r := g.Group("/")
	r.Use(Authorize)
	{
		r.GET("/logout", controller.Logout)
		r.GET("/get-profile", controller.GetProfile)
		r.GET("/list-todo", controller.GetListToDo)
		r.GET("/todo-by-id/:id", controller.GetToDo)
		r.POST("/create-task", controller.CreateToDo)
		r.PUT("/update-task", controller.UpdateToDo)
		r.DELETE("/delete-task/:id", controller.DeleteToDo)
		r.PUT("/update-profile", controller.UpdateProfile)
	}

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
