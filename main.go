package main

import (
	Src "./src"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	guuid "github.com/satori/go.uuid"
)

type boardTable struct {
	gorm.Model
	Mac         string
	Name        string
	Description string
}

type senseDataTable struct {
	gorm.Model
	Mac   string
	Type  string
	Value float64
	Unit  string
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

func main() {
	Src.Db, Src.Err = gorm.Open("sqlite3", "database.db")
	if Src.Err != nil {
		panic("failed to connect database")
	}
	defer Src.Db.Close()
	Src.Db.LogMode(Src.DbLogMode)

	// databasePrepare()
	// Migrate the schema
	Src.Db.AutoMigrate(&boardTable{})
	Src.Db.AutoMigrate(&senseDataTable{})
	Src.Db.AutoMigrate(&Src.BoardSettingsTable{})
	Src.Db.AutoMigrate(&Src.BoardToDoTable{})
	handleHTTP()
}

func handleHTTP() {

	http.HandleFunc("/set_sense_data", sensorDatas)
	http.HandleFunc("/boards", boards)
	http.HandleFunc("/get_board_settings", Src.GetBoardSettings)
	http.HandleFunc("/todo", boardToDo)
	http.HandleFunc("/sensor_types", sensorTypes)
	http.HandleFunc("/sensors_data", sensorDatas)

	fmt.Printf("Starting Server to HANDLE ahome.tech back end\nPort : " + Src.Port + "\nAPI revision " + Src.Version + "\n\n")
	if err := http.ListenAndServe(":"+Src.Port, nil); err != nil {
		log.Fatal(err)
	}
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

		var Response Src.APIHTTPResponseJSONToDo

		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("command_done") != "" {
			Src.Db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", r.URL.Query().Get("mac")).Find(&Response.Entity)
		}
		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("command_done") == "" {
			Src.Db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", false).Limit(2).Find(&Response.Entity)
		}

		Response.API = Src.Version
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

		Src.Db.Create(&Src.BoardToDoTable{Mac: incomingData.Mac,
			Command:     incomingData.Command,
			CommandHash: incomingData.CommandHash,
			CommandDone: incomingData.CommandDone,
			SubCommand:  incomingData.SubCommand})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardToDo Src.BoardToDoTable
		Src.Db.Last(&boardToDo)
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

		var command Src.BoardToDoTable
		Src.Db.Where("command_hash = ?", incomingData.CommandHash).First(&command).Updates(Src.BoardToDoTable{CommandDone: incomingData.CommandDone, CommandStatus: incomingData.CommandStatus})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)

		fmt.Fprintf(w, string(""))

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
			Src.Db.Where("mac = ?", r.URL.Query().Get("mac")).Group("type").Find(&records)
		} else {
			Src.Db.Group("type").Find(&records)
		}

		for _, element := range records {
			Response.Entity = append(Response.Entity, element.Type)
		}

		Response.API = Src.Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func getHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

func sensorDatas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET sensor data")

		var Response apiHTTPResponseJSONSensorDatas
		var sql *gorm.DB

		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") != "" && r.URL.Query().Get("last") != "" {
			sql = Src.Db.Where("mac = ?", r.URL.Query().Get("mac")).
				Where("type = ?", r.URL.Query().Get("type")).
				Last(&Response.Entity)
		} else {

			if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") != "" {
				sql = Src.Db.Where("mac = ?", r.URL.Query().Get("mac")).
					Where("type = ?", r.URL.Query().Get("type"))
			}

			if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("type") == "" {
				sql = Src.Db.Where("mac = ?", r.URL.Query().Get("mac"))
			}

			if r.URL.Query().Get("mac") == "" && r.URL.Query().Get("type") != "" {
				sql = Src.Db.Where("type = ?", r.URL.Query().Get("type"))
			}

			sql.Order("created_at desc").Find(&Response.Entity)
		}

		Response.API = Src.Version
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
		Src.Db.First(&board, "mac = ?", incomingData.Mac)
		fmt.Println(board.Mac)
		if board.Mac == "" {
			Src.Db.Create(&boardTable{Mac: incomingData.Mac})
		}

		Src.Db.Create(&senseDataTable{
			Mac:   incomingData.Mac,
			Type:  incomingData.Valuetype,
			Value: incomingData.Value,
			Unit:  incomingData.Unit})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var senseData senseDataTable
		Src.Db.Last(&senseData)
		addedrecordString, _ := json.Marshal(senseData)

		// fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func boards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, string(""))

	case "POST":
		type incomingDataStructure struct {
			Mac         string `json:"mac"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var incomingData incomingDataStructure

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("content-type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		Src.Db.Create(&boardTable{Mac: incomingData.Mac, Name: incomingData.Name, Description: incomingData.Description})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardData boardTable
		Src.Db.Last(&boardData)
		addedrecordString, _ := json.Marshal(boardData)

		fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	case "GET":

		var Response apiHTTPResponseJSONBoards
		Src.Db.Find(&Response.Entity)
		Response.API = Src.Version
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
		Src.Db.LogMode(Src.DbLogMode)
		Src.Db.First(&board, "mac = ?", incomingData.Mac)
		fmt.Println(board.Mac)
		if board.Mac == "" {
			Src.Db.Create(&boardTable{Mac: incomingData.Mac})
		}
		Src.Db.Create(&senseDataTable{Mac: incomingData.Mac, Type: incomingData.Valuetype, Value: incomingData.Value, Unit: incomingData.Unit})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var senseData senseDataTable
		Src.Db.Last(&senseData)
		addedrecordString, _ := json.Marshal(senseData)

		fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
