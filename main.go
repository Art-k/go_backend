package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	Src "./src"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	
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

	f, err := os.OpenFile("log_go_backend.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// databasePrepare()
	// Migrate the schema
	Src.Db.AutoMigrate(&boardTable{})
	Src.Db.AutoMigrate(&senseDataTable{})
	Src.Db.AutoMigrate(&Src.BoardSettingsTable{})
	Src.Db.AutoMigrate(&Src.BoardToDoTable{})
	Src.Db.AutoMigrate(&Src.UnknownBoards{})
	handleHTTP()
}

func handleHTTP() {

	http.HandleFunc("/set_sense_data", sensorDatas)
	http.HandleFunc("/boards", boards)
	http.HandleFunc("/unknownboards", unknownboards)
	http.HandleFunc("/sensor_types", sensorTypes)
	http.HandleFunc("/sensors_data", sensorDatas)

	http.HandleFunc("/board_settings", Src.GetBoardSettings)
	http.HandleFunc("/todo", Src.BoardToDo)


	fmt.Printf("Starting Server to HANDLE ahome.tech back end\nPort : " + Src.Port + "\nAPI revision " + Src.Version + "\n\n")
	if err := http.ListenAndServe(":"+Src.Port, nil); err != nil {
		log.Fatal(err)
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

			sql.Order("created_at asc").Find(&Response.Entity)
		}

		Response.API = Src.Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
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


func unknownboards(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		type apiHTTPResponseJSONUnknownBoards struct {
			API    string       `json:"api"`
			Total  int          `json:"total"`
			Entity []Src.UnknownBoards `json:"entity"`
		}

		var Response apiHTTPResponseJSONUnknownBoards

		Src.Db.Find(&Response.Entity)
		Response.API = Src.Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
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
