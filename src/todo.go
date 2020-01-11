package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

// BoardToDoTable table
type BoardToDoTable struct {
	gorm.Model
	Mac           string
	Command       string
	SubCommand    string
	CommandHash   string
	CommandDone   bool
	CommandStatus string
}

type apiHTTPResponseJSONToDo struct {
	API    string           `json:"api"`
	Total  int              `json:"total"`
	Entity []BoardToDoTable `json:"entity"`
}

// BoardToDo requests handler
func BoardToDo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":

		log.Println("/todo OPTIONS received")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		n, _ := fmt.Fprintf(w, string(""))
		fmt.Println(n)

		log.Println("/todo OPTIONS answered\n\n")

	case "GET":

		log.Println("/todo GET received")
		var Response apiHTTPResponseJSONToDo

		log.Println("/todo check GET parameters")
		log.Println("/todo mac='" + r.URL.Query().Get("mac") + "' and command_done='" + r.URL.Query().Get("command_done") + "'")

		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("command_done") != "" {

			Db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", r.URL.Query().Get("command_done")).Order("created_at desc").Find(&Response.Entity)

		}

		if r.URL.Query().Get("mac") != "" && r.URL.Query().Get("command_done") == "" {
			if r.URL.Query().Get("limit") != "" {
				Db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", false).Order("created_at asc").Limit(r.URL.Query().Get("limit")).Find(&Response.Entity)
			} else {
				Db.Where("mac = ?", r.URL.Query().Get("mac")).Where("command_done = ?", false).Order("created_at desc").Find(&Response.Entity)
			}

		}

		Response.API = Version
		Response.Total = len(Response.Entity)

		addedRecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)

		log.Println("/todo GET answered\n\n")

	case "POST":

		log.Println("/todo POST received")

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

		incomingData.CommandHash = GetHash()
		incomingData.CommandDone = false

		fmt.Println(incomingData)

		Db.Create(&BoardToDoTable{
			Mac:         incomingData.Mac,
			Command:     incomingData.Command,
			CommandHash: incomingData.CommandHash,
			CommandDone: incomingData.CommandDone,
			SubCommand:  incomingData.SubCommand,
		})

		log.Println("/todo record added")
		log.Println("/todo \t Mac :\t" + incomingData.Mac)
		log.Println("/todo \t Command :\t" + incomingData.Command)
		log.Println("/todo \t CommandHash :\t" + incomingData.CommandHash)
		log.Println("/todo \t CommandDone :\t" + strconv.FormatBool(incomingData.CommandDone))
		log.Println("/todo \t SubCommand :\t" + incomingData.SubCommand)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardToDo BoardToDoTable
		Db.Last(&boardToDo)
		addedrecordString, _ := json.Marshal(boardToDo)
		fmt.Fprintf(w, string(addedrecordString))

		log.Println("/todo POST done\n\n")

	case "PATCH":

		log.Println("/todo PATCH received")

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

		log.Println("/todo \t CommandHash : " + incomingData.CommandHash)
		log.Println("/todo \t CommandDone : " + strconv.FormatBool(incomingData.CommandDone))
		log.Println("/todo \t CommandStatus : " + incomingData.CommandStatus)

		var command BoardToDoTable
		Db.Where("command_hash = ?", incomingData.CommandHash).First(&command).Updates(BoardToDoTable{CommandDone: incomingData.CommandDone, CommandStatus: incomingData.CommandStatus})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)

		fmt.Fprintf(w, string(""))

		log.Println("/todo PATCH done\n\n")

	default:
		fmt.Fprintf(w, "Sorry, only OPTIONS,GET,POST,PATCH methods are supported. '"+r.Method+"' received")
		log.Println("/todo PATCH done\n\n")
	}
}
