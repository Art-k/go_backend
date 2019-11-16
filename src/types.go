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
