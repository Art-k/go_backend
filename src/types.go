package src

import "github.com/jinzhu/gorm"

var Db *gorm.DB
var Err error

const Version = "0.2.1"
const DbLogMode = true
const Port = "55555"
