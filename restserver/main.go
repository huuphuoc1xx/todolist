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

	g.POST("/user/login", controller.Login)
	g.POST("/user/register", controller.Register)
	g.GET("/user/logout", Authorize, controller.Logout)

	r := g.Group("/user")
	r.Use(Authorize)
	r.GET("/", controller.GetProfile)
	r.PUT("/", controller.UpdateProfile)

	r2 := g.Group("/todo/")
	r2.Use(Authorize)
	r2.GET("/bytag", controller.GetListToDo)
	r2.GET("/byid/:id", controller.GetToDo)
	r2.POST("/", controller.CreateToDo)
	r2.PUT("/", controller.UpdateToDo)
	r2.DELETE("/", controller.DeleteToDo)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
