package router

import (
	"todolist/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetListToDo(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pagenumber"})
		return
	}

	tag:=ctx.Query("tag")
	if len(tag)==0{
		tag="*"
	}
	util.GetTaskByPage(db)
}
