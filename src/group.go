package src

import (
	"encoding/json"
	// "encoding/json"
	// "fmt"
	// "log"
	"fmt"
	"log"
	"net/http"
	// "strconv"
)

// Groups list of groups
func Groups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
	case "GET":

		type APIHTTPResponseJSONGroups struct {
			API    string   `json:"api"`
			Total  int      `json:"total"`
			Entity []SGroup `json:"entity"`
		}

		var Response APIHTTPResponseJSONGroups

		Db.Find(&Response.Entity)

		Response.API = Version
		Response.Total = len(Response.Entity)
		addedRecordString, _ := json.Marshal(Response)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only OPTIONS,GET,POST,PATCH methods are supported. '"+r.Method+"' received")
		log.Println(n)
		log.Println("/todo PATCH done\n\n")
	}
}

// GroupCRUD crud for one group
func GroupCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		log.Println("/board_settings OPTIONS received")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		n, _ := fmt.Fprintf(w, string(""))
		fmt.Println(n)

		log.Println("/board_settings OPTIONS answered")

	case "POST":

		fmt.Println("POST group data")

		var incomingData GroupPost

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Check if the group exists
		var group SGroup
		// db.LogMode(dbLogMode)
		Db.First(&group, "name = ?", incomingData.Name)
		fmt.Println(group.Name)
		if group.Name == "" {
			Db.Create(&SGroup{Name: incomingData.Name})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		Db.First(&group, "name = ?", incomingData.Name)
		addedRecordString, _ := json.Marshal(group)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)

	case "PATCH":
	case "DELETE":
	case "GET":
	default:
		n, _ := fmt.Fprintf(w, "Sorry, only OPTIONS,GET,POST,PATCH methods are supported. '"+r.Method+"' received")
		log.Println(n)
		log.Println("/todo PATCH done\n\n")
	}
}
