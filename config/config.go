package config
import (
	"time"
)
type ToDo struct {
	Username string `json:"username"`
	Id int `json:"id"`
	Title string `json:"Title"`
	Tag string `json:"Tag"`
	Description string `json:"Description"`
	InsertTime time.Time `json:"InsertTime"` 
}