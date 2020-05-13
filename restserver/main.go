package main
import (
	"github.com/gin-gonic/gin""
)
func main()  {
	g:=gin.Default();
	g.GET("/list-by-page",getListPaginate)

	g.GET("/search-by-page",getListByTag)

	g.POST("/new-task",createTask)

	g.PUT("/")
}