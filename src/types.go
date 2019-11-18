package src

import "github.com/jinzhu/gorm"

// Db database
var Db *gorm.DB

// Err error
var Err error

// Version of api
const Version = "0.2.1"

// DbLogMode log mode for database
const DbLogMode = true

// Port application use this port to get requests
const Port = "55555"

type SenseDataTable struct {
	gorm.Model
	Mac   string
	Type  string
	Value float64
	Unit  string
}

type APIHTTPResponseJSONSensorDatas struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []SenseDataTable `json:"entity"`
}
