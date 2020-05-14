package controller

import (
	"database/sql"
	"log"
	"time"
	"todolist/config"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Username = "gsMg5DbNgQ"
	Password = "32F8Gb0lfr"
	url      = "remotemysql.com"
	port     = "3306"
	database = Username
)

var db = initDB()

func initDB() *sql.DB {
	db, err := sql.Open("mysql", Username+":"+Password+"@tcp"+"("+url+":"+port+")/"+database)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func getTaskByPage(tag string, page int) []config.ToDo {
	rows, err := db.Query("Select * from todolist where Tag=? limit 20 offset ?", tag, page-1)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	todos := make([]config.ToDo, 0)
	for rows.Next() {
		var todo config.ToDo
		var insertTime string
		if err := rows.Scan(&todo.Username, &todo.Id, &todo.Title, &todo.Tag, &todo.Description, &insertTime); err != nil {
			log.Fatal(err)
		}

		formatTime := "2006-01-02 15:04:05"
		todo.InsertTime, _ = time.Parse(formatTime, insertTime)
		todos = append(todos, todo)
	}

	return todos
}

func createTask(todo config.ToDo) int {
	res, err := db.Exec("Insert into todolist(Username,Title,Tag,Description) value (?,?,?,?)", todo.Username, todo.Title, todo.Tag, todo.Description)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		return 0
	}
	id, err := res.LastInsertId()
	return int(id)
}

func updateTask(todo config.ToDo) int {
	res, err := db.Exec("Update todolist Set Username=?, Title=?,Tag=?,Description=?where ID=?",
		todo.Username, todo.Title, todo.Tag, todo.Description, todo.Id)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		return 0
	}
	id, err := res.LastInsertId()
	return int(id)
}

func getTaskById(id int) []config.ToDo {
	row := db.QueryRow("Select * from todolist where ID=?", id)

	var todo config.ToDo
	var insertTime string
	if err := row.Scan(&todo.Username, &todo.Id, &todo.Title, &todo.Tag, &todo.Description, &insertTime); err != nil {
		return nil
	}

	formatTime := "2006-01-02 15:04:05"
	todo.InsertTime, _ = time.Parse(formatTime, insertTime)
	res := make([]config.ToDo, 0)
	res = append(res, todo)
	return res
}

func deleteTask(id int) int {
	res, err := db.Exec("Delete todolist where ID=?", id)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		return 0
	}
	return id
}
