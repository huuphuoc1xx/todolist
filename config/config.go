package config

import (
	"database/sql"
	"time"
)

type ToDo struct {
	Username    string    `json:"username"`
	Id          int       `json:"id"`
	Title       string    `json:"Title"`
	Tag         string    `json:"Tag"`
	Description string    `json:"Description"`
	InsertTime  time.Time `json:"InsertTime"`
}

type DBInfo struct {
	username string
	password string
	url      string
	port     string
	database string
}

func (db *DBInfo) SetTodoDB() {
	db.username = "gsMg5DbNgQ"
	db.password = "0OR7XKmLK5"
	db.url = "remotemysql.com"
	db.port = "3306"
	db.database = "gsMg5DbNgQ"
}

func (db *DBInfo) SetProfileDB() {
	db.username = "gsMg5DbNgQ"
	db.password = "0OR7XKmLK5"
	db.url = "remotemysql.com"
	db.port = "3306"
	db.database = "gsMg5DbNgQ"
}

func (db *DBInfo) GetDB() (*sql.DB, error) {
	return sql.Open("mysql", db.username+":"+db.password+"@tcp"+"("+db.url+":"+db.port+")/"+db.database)
}
