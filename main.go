package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	guuid "github.com/satori/go.uuid"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const version = "0.2.1"
const dbLogMode = true
const port = "55555"

var db *gorm.DB
var err error

type boardTable struct {
	gorm.Model
	Mac         string
	Name        string
	Description string
}

type boardSettingsTable struct {
	gorm.Model
	Mac                  string
	SensorType           string
	Sense                string
	Pin                  string
	Interval             int64
	Delta                int64
	Default              string
	AdditionalParameters string
}

type senseDataTable struct {
	gorm.Model
	Mac   string
	Type  string
	Value float64
	Unit  string
}

type boardToDoTable struct {
	gorm.Model
	Mac           string
	Command       string
	SubCommand    string
	CommandHash   string
	CommandDone   bool
	CommandStatus string
}

type apiHTTPResponseJSONBoards struct {
	API    string       `json:"api"`
	Total  int          `json:"total"`
	Entity []boardTable `json:"entity"`
}

type apiHTTPResponseJSONSensorDatas struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []senseDataTable `json:"entity"`
}

type apiHTTPResponseJSONSensorTypes struct {
	API    string   `json:"api"`
	Total  int      `json:"total"`
	Entity []string `json:"entity"`
}

type apiHTTPResponseJSONBoardSetttings struct {
	API    string               `json:"api"`
	Total  int                  `json:"total"`
	Entity []boardSettingsTable `json:"entity"`
}

type apiHTTPResponseJSONToDo struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []boardToDoTable `json:"entity"`
}

func main() {
	db, err = gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(dbLogMode)

	// databasePrepare()
	// Migrate the schema
	db.AutoMigrate(&boardTable{})
	db.AutoMigrate(&senseDataTable{})
	db.AutoMigrate(&boardSettingsTable{})
	db.AutoMigrate(&boardToDoTable{})
	handleHTTP()
}

func handleHTTP() {

	http.HandleFunc("/set_sense_data", sensorDatas)
	http.HandleFunc("/get_boards", getBoards)
	http.HandleFunc("/get_board_settings", getBoardSettings)
	http.HandleFunc("/todo", boardToDo)
	http.HandleFunc("/sensor_types", sensorTypes)
	http.HandleFunc("/sensors_data", sensorDatas)

	fmt.Printf("Starting Server to HANDLE ahome.tech back end\nPort : " + port + "\nAPI revision " + version + "\n\n")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

func boardToDo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, string(""))
	case "GET":

		var Response apiHTTPResponseJSONToDo

		if r.URL.Query().Get("mac") != "" {
			db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", false).Find(&Response.Entity)
		}

		Response.API = version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	case "POST":
		fmt.Println("POST command")
		type incomingDataStructure struct {
			Mac           string
			Command       string
			SubCommand    string
			CommandHash   string
			CommandDone   bool
			CommandStatus string
		}
		var incomingData incomingDataStructure

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		incomingData.CommandHash = getHash()
		incomingData.CommandDone = false

		fmt.Println(incomingData)

		db.Create(&boardToDoTable{Mac: incomingData.Mac,
			Command:     incomingData.Command,
			CommandHash: incomingData.CommandHash,
			CommandDone: incomingData.CommandDone,
			SubCommand:  incomingData.SubCommand})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardToDo boardToDoTable
		db.Last(&boardToDo)
		addedrecordString, _ := json.Marshal(boardToDo)

		// fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	case "PATCH":

		fmt.Println("PATCH command")
		type incomingDataStructure struct {
			Mac           string
			CommandHash   string
			CommandDone   bool
			CommandStatus string
		}
		var incomingData incomingDataStructure

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(incomingData)

		var command boardToDoTable
		db.Where("command_hash = ?", incomingData.CommandHash).First(&command).Updates(boardToDoTable{CommandDone: incomingData.CommandDone, CommandStatus: incomingData.CommandStatus})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)

		fmt.Fprintf(w, string(""))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func getBoardSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		var Response apiHTTPResponseJSONBoardSetttings

		if r.URL.Query().Get("mac") != "" {
			db.Where("mac = ?", r.URL.Query().Get("mac")).Find(&Response.Entity)
		}

		Response.API = version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func sensorTypes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		type responseJSON struct {
			API    string   `json:"api"`
			Total  int      `json:"total"`
			Entity []string `json:"entity"`
		}
		var Response responseJSON
		var records []senseDataTable
		if r.URL.Query().Get("mac") != "" {
			db.Where("mac = ?", r.URL.Query().Get("mac")).Group("type").Find(&records)
		} else {
			db.Group("type").Find(&records)
		}

		for _, element := range records {
			Response.Entity = append(Response.Entity, element.Type)
		}

		Response.API = version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func sensorDatas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET sensor data")

		var Response apiHTTPResponseJSONSensorDatas
		var sql *gorm.DB

		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") != "" && r.URL.Query().Get("last") != "" {
			sql = db.Where("mac = ?", r.URL.Query().Get("mac")).
				Where("type = ?", r.URL.Query().Get("type")).
				Last(&Response.Entity)
		} else {

			if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") != "" {
				sql = db.Where("mac = ?", r.URL.Query().Get("mac")).
					Where("type = ?", r.URL.Query().Get("type"))
			}

			if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") == "" {
				sql = db.Where("mac = ?", r.URL.Query().Get("mac"))
			}

			if r.URL.Query().Get("mac") == "" && r.URL.Query().Get("type") != "" {
				sql = db.Where("type = ?", r.URL.Query().Get("type"))
			}

			sql.Order("created_at desc").Find(&Response.Entity)
		}

		Response.API = version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	case "POST":
		fmt.Println("POST sensor data")
		type incomingDataStructure struct {
			Mac       string  `json:"mac"`
			Valuetype string  `json:"valuetype"`
			Value     float64 `json:"value"`
			Unit      string  `json:"unit"`
		}
		var incomingData incomingDataStructure

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Check if the board exists
		var board boardTable
		// db.LogMode(dbLogMode)
		db.First(&board, "mac = ?", incomingData.Mac)
		fmt.Println(board.Mac)
		if board.Mac == "" {
			db.Create(&boardTable{Mac: incomingData.Mac})
		}

		db.Create(&senseDataTable{Mac: incomingData.Mac, Type: incomingData.Valuetype, Value: incomingData.Value, Unit: incomingData.Unit})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var senseData senseDataTable
		db.Last(&senseData)
		addedrecordString, _ := json.Marshal(senseData)

		// fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func getBoards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		var Response apiHTTPResponseJSONBoards
		db.Find(&Response.Entity)
		Response.API = version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func setSensorData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		type incomingDataStructure struct {
			Mac       string  `json:"mac"`
			Valuetype string  `json:"valuetype"`
			Value     float64 `json:"value"`
			Unit      string  `json:"unit"`
		}
		var incomingData incomingDataStructure

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Check if the board exists
		var board boardTable
		db.LogMode(dbLogMode)
		db.First(&board, "mac = ?", incomingData.Mac)
		fmt.Println(board.Mac)
		if board.Mac == "" {
			db.Create(&boardTable{Mac: incomingData.Mac})
		}
		db.Create(&senseDataTable{Mac: incomingData.Mac, Type: incomingData.Valuetype, Value: incomingData.Value, Unit: incomingData.Unit})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var senseData senseDataTable
		db.Last(&senseData)
		addedrecordString, _ := json.Marshal(senseData)

		fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
