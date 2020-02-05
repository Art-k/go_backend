package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type httpIncomingLog struct {
	SessionId string `json:"s"`
	Mac       uint   `json:"mac"`
	Log       string `json:"str"`
}

type apiHTTPResponseJSONBoardLogs struct {
	API    string     `json:"api"`
	Total  int        `json:"total"`
	Entity []BoardLog `json:"entity"`
}

func ActionLog(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		OptionsAnswer(w, r)
		return
	case "POST":

		var incomingData httpIncomingLog

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(incomingData)

		Db.Create(&BoardLog{
			Mac:       incomingData.Mac,
			SessionId: incomingData.SessionId,
			Log:       incomingData.Log,
		})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, "")
		fmt.Println(n)
		return

	case "PATCH":

	case "GET":
	case "DELETE":
	default:
	}
}

func ActionLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("/board_log GET received")

		var Response apiHTTPResponseJSONBoardLogs

		if r.URL.Query().Get("mac") != "" {
			Db.Where("mac = ?", r.URL.Query().Get("mac")).Find(&Response.Entity)
		}

		Response.API = Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedrecordString))
		fmt.Println(n)

		log.Println("/board_log GET answered\n\n")
		return

	default:

	}
}
