package src

import "github.com/jinzhu/gorm"

// Db database
var Db *gorm.DB

// Err error
var Err error

// Version of api
const Version = "0.2.1"

// DbLogMode log mode for database
const DbLogMode = false

// Port application use this port to get requests
const Port = "55555"

// SenseDataTable collect sense data to table
type SenseDataTable struct {
	gorm.Model
	Mac   string
	Type  string
	Value float64
	Unit  string
}

// APIHTTPResponseJSONSensorDatas respons
type APIHTTPResponseJSONSensorDatas struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []SenseDataTable `json:"entity"`
}

// IncomingDataStructure structure
type IncomingDataStructure struct {
	Mac       string  `json:"mac"`
	Valuetype string  `json:"valuetype"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

// DeviceState state dynamic
type DeviceState struct {
	gorm.Model
	ByMac    string `json:"by_mac"`
	NewState string `json:"new_state"`
}

// SGroup sensor groups to show in one screen or chart
type SGroup struct {
	gorm.Model
	Name    string         `json:"name;unique_index"`
	Sensors []SensorsGroup `json:"sensors"`
}

type GroupPost struct {
	Name string `json:"name"`
}

type GroupPatchDelete struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SensorsGroup groups, all sensors and boards in group
type SensorsGroup struct {
	gorm.Model
	GroupID    uint   `json:"group_id"`
	MacID      uint   `json:"mac_id"`
	SensorType string `json:"sensor_type"`
}

// SensorsGroup groups, all sensors and boards in group
type BoardLog struct {
	gorm.Model
	SessionId int    `json:"s"`
	Mac       uint   `json:"mac"`
	Log       string `json:"str"`
}
