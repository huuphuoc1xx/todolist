package controller

import (
	"database/sql"
	"log"
	"time"
	"todolist/config"
)

const (
	Username = "gsMg5DbNgQ"
	Password = "9gkFJA1OWf"
	url      = "remotemysql.com"
	port     = "3306"
	database = Username
)

var db = initDB()

func initDB() *sql.DB {
	db, err := sql.Open("mysql", Username+":"+Password+"@tcp"+"("+url+":"+port+")/"+database)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func GetTaskByPage(tag string, page int) []config.ToDo {
	rows, err := db.Query("Select * from todolist where Tag=? limit 20 offset ?", tag, page-1)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	todos := make([]config.ToDo, 0)
	for rows.Next() {
		var todo config.ToDo
		var insertTime string
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Tag, &todo.Description, &insertTime); err != nil {
			log.Fatal(err)
		}

		formatTime := "2006-01-02 15:04:05"
		todo.InsertTime, _ = time.Parse(formatTime, insertTime)
		todos = append(todos, todo)
	}

	log.Fatal(todos)

	return todos
}
