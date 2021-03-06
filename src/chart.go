package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// GetChartData http handler

// GetChartsData the function to return charts data several boards + several
func GetChartsData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		log.Println("/charts GET received")
		var startDate = r.URL.Query().Get("start")
		var endDate = r.URL.Query().Get("end")
		var group = r.URL.Query().Get("group")

		fmt.Println("Start Date", startDate)
		fmt.Println("End Date", endDate)
		fmt.Println("Group by", group)

		type senseData struct {
			mac   string
			sense string
		}
		var incomingData []senseData

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			log.Println("Wrong Input")
		}

		type rec []float32
		type responseData struct {
			datas  []string
			values []rec
			header []string
		}

		// for _, elem := range incomingData{
		// 	Data =
		// }

	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported. '"+r.Method+"' received")
		log.Println("/charts GET done\n")
	}

}

// GetChartData simple chart for one parameter
func GetChartData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":

		log.Println("/chart GET received\n")

		var MAC = r.URL.Query().Get("mac")
		var startDate = r.URL.Query().Get("start")
		var endDate = r.URL.Query().Get("end")
		var senseType = r.URL.Query().Get("type")
		var Compact = r.URL.Query().Get("compact")

		log.Println("/chart MAC : " + MAC)
		if MAC != "" {
			log.Println("/chart senseType : " + senseType)
			if senseType == "" {
				senseType = getDefaultSenseType()
			}

			log.Println("/chart startDate : " + startDate)
			if startDate == "" {
				startDate = getDefaultStartDate()
			}

			log.Println("/chart endDate : " + endDate)
			if endDate == "" {
				endDate = getDefaultEndDate()
			}

			var Response APIHTTPResponseJSONSensorDatas
			Db.Where("mac = ?", MAC).Where("type = ?", senseType).Where("created_at <= ?", endDate).Where("created_at >= ?", startDate).Find(&Response.Entity)

			if Compact == "1" {
				var PrevValue float64
				PrevValue = -99999999
				var PrevSense string
				PrevSense = ""
				var NewEntity []SenseDataTable
				for _, element := range Response.Entity {
					if PrevSense == "" {
						PrevSense = element.Type
						if PrevValue == -99999999 {
							PrevValue = element.Value
						}
						continue
					}
					if !(PrevSense == element.Type && PrevValue == element.Value) {
						NewEntity = append(NewEntity, element)
						PrevSense = element.Type
						PrevValue = element.Value
					}
				}

				Response.Entity = NewEntity
			}

			Response.API = Version
			Response.Total = len(Response.Entity)

			addedrecordString, _ := json.Marshal(Response)

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("content-type", "application/json")
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

func removeSenseDataTable(slice []SenseDataTable, s int) []SenseDataTable {
	return append(slice[:s], slice[s+1:]...)
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

	// type responseJSON struct {
	// 	API    string   `json:"api"`
	// 	Total  int      `json:"total"`
	// 	Entity []string `json:"entity"`
	// }

	// // var Response responseJSON
	// var records []SenseDataTable
	// Db.Group("type").Find(&records)

	// var Result string
	// for _, element := range records {
	// 	Result = element.Type
	// 	break
	// }

	return "temperature"
}
