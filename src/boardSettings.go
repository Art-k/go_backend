package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

// UnknownBoards list of the boards which does't has settings
type UnknownBoards struct {
	gorm.Model
	Mac string
}

// BoardSettingsTable table
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

// GetBoardSettings request handler
func GetBoardSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "OPTIONS":

		log.Println("/board_settings OPTIONS received")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, string(""))

		log.Println("/board_settings OPTIONS answered\n\n")

	case "DELETE":

		log.Println("/board_settings DELETE received")

		// var boardSettings BoardSettingsTable
		if r.URL.Query().Get("id") != "" {
			log.Println("/board_settings record with id=" + r.URL.Query().Get("id") + " will be deleted")
			Db.Where("id = ?", r.URL.Query().Get("id")).Delete(&BoardSettingsTable{})
		}

		log.Println("/board_settings DELETE done\n\n")

	case "POST":

		log.Println("/board_settings POST received")

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
			log.Println("/board_settings ERROR in incoming JSON")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("content-type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		Db.Create(&BoardSettingsTable{Mac: incomingData.Mac,
			SensorType:           incomingData.SensorType,
			Sense:                incomingData.Sense,
			Pin:                  incomingData.Pin,
			Interval:             incomingData.Interval,
			Delta:                incomingData.Delta,
			Default:              incomingData.Default,
			AdditionalParameters: incomingData.AdditionalParameters})

		log.Println("/board_settings record added")
		log.Println("/board_settings \t Mac :\t" + incomingData.Mac)
		log.Println("/board_settings \t SensorType :\t" + incomingData.SensorType)
		log.Println("/board_settings \t Sense :\t" + incomingData.Sense)
		log.Println("/board_settings \t Pin :\t" + incomingData.Pin)
		log.Println("/board_settings \t Interval :\t" + strconv.FormatInt(incomingData.Interval, 10))
		log.Println("/board_settings \t Delta :\t" + strconv.FormatInt(incomingData.Delta, 10))
		log.Println("/board_settings \t Default :\t" + incomingData.Default)
		log.Println("/board_settings \t AdditionalParameters :\t" + incomingData.AdditionalParameters)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusCreated)

		var boardSettingsData BoardSettingsTable
		Db.Last(&boardSettingsData)
		addedrecordString, _ := json.Marshal(boardSettingsData)

		fmt.Println(addedrecordString)
		fmt.Fprintf(w, string(addedrecordString))

		log.Println("/board_settings POST done\n\n")

	case "GET":

		log.Println("/board_settings GET received")

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

		log.Println("/board_settings GET answered\n\n")

		if Response.Total == 0 {
			Db.Create(&UnknownBoards{Mac: r.URL.Query().Get("mac")})
			log.Println("/board_settings GET MAC '" + r.URL.Query().Get("mac") + "' is unknown, added to list of unknown boards \n\n")
		}

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
