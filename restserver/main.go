package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	g := gin.Default()
	g.GET("/list-by-page")

	g.GET("/search-by-page", getListByTag)

	g.POST("/new-task", createTask)

	g.PUT("/")
}
