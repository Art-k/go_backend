package src

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

type RuleBySensor struct {
	gorm.Model
	SensorMac   string
	SensorType  string
	ActionMac   string
	Condition   string
	DoIFTrue    string
	DoIFFalse   string
	RepeatEvery int64
	Expires     int64
	Active      bool
}

type RuleByTimer struct {
	gorm.Model
	ActionMac   string
	DoIFTrue    string
	DoIFFalse   string
	StartsAt    int64
	RepeatEvery int64
	Expires     int64
	Active      bool
}

type Rule struct {
	RuleType     string //timer, sensor
	RuleByTimer  RuleByTimer
	RuleBySensor RuleBySensor
}

func ActionRulePOST(w http.ResponseWriter, r *http.Request) {

	type ruleBySensor struct {
		SensorMac   string
		SensorType  string
		ActionMac   string
		Condition   string
		DoIFTrue    string
		DoIFFalse   string
		RepeatEvery int64
		Expires     int64
		Active      bool
	}

	type ruleByTimer struct {
		ActionMac   string
		DoIFTrue    string
		DoIFFalse   string
		StartsAt    int64
		RepeatEvery int64
		Expires     int64
		Active      bool
	}

	type incomingDataStructure struct {
		RuleType     string //timer, sensor
		RuleByTimer  RuleByTimer
		RuleBySensor RuleBySensor
	}

	var incomingData incomingDataStructure
	err := json.NewDecoder(r.Body).Decode(&incomingData)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if incomingData.RuleType == "sensor" {
		Db.Create(&RuleBySensor{
			SensorMac:   incomingData.RuleBySensor.SensorMac,
			SensorType:  incomingData.RuleBySensor.SensorType,
			ActionMac:   incomingData.RuleBySensor.ActionMac,
			Condition:   incomingData.RuleBySensor.Condition,
			DoIFTrue:    incomingData.RuleBySensor.DoIFTrue,
			DoIFFalse:   incomingData.RuleBySensor.DoIFFalse,
			RepeatEvery: incomingData.RuleBySensor.RepeatEvery,
			Expires:     incomingData.RuleBySensor.Expires,
			Active:      incomingData.RuleBySensor.Active,
		})

		var rule ruleBySensor
		Db.Last(&rule)
		addedRecordString, _ := json.Marshal(rule)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)
		return
	}

	if incomingData.RuleType == "timer" {
		Db.Create(&RuleByTimer{
			ActionMac:   incomingData.RuleByTimer.ActionMac,
			DoIFTrue:    incomingData.RuleByTimer.DoIFTrue,
			DoIFFalse:   incomingData.RuleByTimer.DoIFFalse,
			StartsAt:    incomingData.RuleByTimer.StartsAt,
			RepeatEvery: incomingData.RuleByTimer.RepeatEvery,
			Expires:     incomingData.RuleByTimer.Expires,
			Active:      incomingData.RuleByTimer.Active,
		})

		var rule ruleByTimer
		Db.Last(&rule)
		addedRecordString, _ := json.Marshal(rule)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)
		return
	}

}

func ActionRule(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		OptionsAnswer(w, r)
		return
	case "POST":
		ActionRulePOST(w, r)
		return

	case "PATCH":
	case "GET":
	case "DELETE":
		id := r.URL.Query().Get("id")
		Type := r.URL.Query().Get("type")

		if Type == "timer" {
			var rule RuleByTimer
			Db.Where("id = ?", id).Delete(&rule)

			addedRecordString, _ := json.Marshal(rule)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedRecordString))
			fmt.Println(n)
			return
		}

		if Type == "sensor" {
			var rule RuleBySensor
			Db.Where("id = ?", id).Delete(&rule)

			addedRecordString, _ := json.Marshal(rule)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedRecordString))
			fmt.Println(n)
			return
		}

	default:
	}
}

func ActionRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		type apiHTTPResponseJSONRule struct {
			API    string        `json:"api"`
			Total  int           `json:"total"`
			Entity []RuleByTimer `json:"entity"`
		}

		//active := r.URL.Query().Get("active")
		//expires := r.URL.Query().Get("show_expires")

		ct := int64(time.Now().Unix())
		var response apiHTTPResponseJSONRule
		response.API = Version
		Db.Where("active = ?", true).Where("expires >= ?", ct).Find(&response.Entity)
		//Db.Where("active = ?", true).Find(&response.Entity)

		response.Total = len(response.Entity)
		addedRecordString, _ := json.Marshal(response)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)

	default:
	}
}

func CheckIfWeHaveARule(t time.Time) {
	//fmt.Println("Tik")
	ct := int64(time.Now().Unix())
	log.Println(ct)
	var activeRulesByTimer []RuleByTimer
	Db.Where("active = ?", true).Where("expires >= ?", ct).Find(&activeRulesByTimer)

	for _, rule := range activeRulesByTimer {
		if (ct-rule.StartsAt)%rule.RepeatEvery == 0 {

			Db.Where("mac = ?", rule.ActionMac).
				Where("command_done = ?", false).
				Where("command_sent = ?", false).
				Where("command_status = ?", "").Delete(&BoardToDoTable{})

			Db.Create(&BoardToDoTable{
				Mac:         rule.ActionMac,
				Command:     "RELAY",
				CommandHash: GetHash(),
				CommandDone: false,
				SubCommand:  rule.DoIFTrue,
			})

			postTelegrammMessage("Send Command : " + rule.DoIFTrue)

			fmt.Println("=======================================")
			fmt.Println("=======================================")
			fmt.Println("=======================================")

		}
	}
}
