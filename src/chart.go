package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// GetChartData http handler
func GetChartData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":

		log.Println("/chart GET received\n")

		var MAC = r.URL.Query().Get("mac")
		var startDate = r.URL.Query().Get("start")
		var endDate = r.URL.Query().Get("end")
		// var senseType = r.URL.Query().Get("type")

		log.Println("/chart MAC : " + MAC)
		if MAC != "" {
			// log.Println("/chart senseType : " + senseType)
			// if senseType == "" {
			// 	senseType = getDefaultSenseType()
			// }

			log.Println("/chart startDate : " + startDate)
			if startDate == "" {
				startDate = getDefaultStartDate()
			}

			log.Println("/chart endDate : " + endDate)
			if endDate == "" {
				endDate = getDefaultEndDate()
			}

			var sD []SenseDataTable
			Db.Where("mac = ?", MAC).Where("created_at <= ?", endDate).Where("created_at >= ?", startDate).Group("type").Find(&sD)

			var sData []SenseDataTable
			Db.Where("mac = ?", MAC).Where("created_at <= ?", endDate).Where("created_at >= ?", startDate).Find(&sData)

			var Response APIHTTPResponseJSONSensorDatas

			Response.API = Version
			Response.Total = len(Response.Entity)

			addedrecordString, _ := json.Marshal(Response)

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(addedrecordString))

			log.Println("/todo GET answered\n\n")

		} else {
			fmt.Fprintf(w, "There is no MAC, ")
			log.Println("/chart GET done\n")
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported. '"+r.Method+"' received")
		log.Println("/chart GET done\n")
	}
}

func getDefaultStartDate() string {
	var yesterdayTime = time.Now().Add(-24 * time.Hour)
	return yesterdayTime.Format("2006-01-02")
}

func getDefaultEndDate() string {
	var currentTime = time.Now()
	return currentTime.Format("2006-01-02")
}

func getDefaultSenseType() string {
	// var senseData SenseDataTable

	type responseJSON struct {
		API    string   `json:"api"`
		Total  int      `json:"total"`
		Entity []string `json:"entity"`
	}

	// var Response responseJSON
	var records []SenseDataTable
	Db.Group("type").Find(&records)

	var Result string
	for _, element := range records {
		Result = element.Type
		break
	}

	return Result
}
