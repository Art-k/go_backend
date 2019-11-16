package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

type BoardSettingsTable struct {
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

type apiHTTPResponseJSONBoardSetttings struct {
	API    string               `json:"api"`
	Total  int                  `json:"total"`
	Entity []BoardSettingsTable `json:"entity"`
}

// GetBoardSettings
func GetBoardSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, string(""))

	// TODO: This section should be finished i mean POST for Board settings
	case "POST":
		type incomingDataStructure struct {
			Mac                  string
			SensorType           string
			Sense                string
			Pin                  string
			Interval             int64
			Delta                int64
			Default              string
			AdditionalParameters string
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

		Db.Create(&BoardSettingsTable{Mac: incomingData.Mac, Name: incomingData.Name, Description: incomingData.Description})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardSettingsData BoardSettingsTable
		Db.Last(&boardSettingsData)
		addedrecordString, _ := json.Marshal(boardSettingsData)

		fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

	case "GET":

		var Response apiHTTPResponseJSONBoardSetttings

		if r.URL.Query().Get("mac") != "" {
			Db.Where("mac = ?", r.URL.Query().Get("mac")).Find(&Response.Entity)
		}

		Response.API = Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
