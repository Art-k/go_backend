package src

import "github.com/jinzhu/gorm"

/*
	"BoardToDoTable"
*/
type BoardToDoTable struct {
	gorm.Model
	Mac           string
	Command       string
	SubCommand    string
	CommandHash   string
	CommandDone   bool
	CommandStatus string
}

/*
	"APIHTTPResponseJSONToDo"
*/
type APIHTTPResponseJSONToDo struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []BoardToDoTable `json:"entity"`
}
