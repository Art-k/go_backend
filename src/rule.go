package src

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Rule struct {
	gorm.Model
	SensorMac  string
	SensorType string
	ActionMac  string
	Condition  string
	DoIFTrue   string
	DoIFFalse  string
}

func ActionRulePOST(w http.ResponseWriter, r *http.Request) {

	type incomingDataStructure struct {
		SensorMac  string
		SensorType string
		ActionMac  string
		Condition  string
		DoIFTrue   string
		DoIFFalse  string
	}

	var incomingData incomingDataStructure
	err := json.NewDecoder(r.Body).Decode(&incomingData)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Db.Create(&Rule{
		SensorMac:  incomingData.SensorMac,
		SensorType: incomingData.SensorType,
		ActionMac:  incomingData.ActionMac,
		Condition:  incomingData.Condition,
		DoIFTrue:   incomingData.DoIFTrue,
		DoIFFalse:  incomingData.DoIFFalse,
	})

	var rule Rule
	Db.Last(&rule)
	addedRecordString, _ := json.Marshal(rule)
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, string(addedRecordString))
	fmt.Println(n)
	return
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
	default:
	}
}

func ActionRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	default:
	}
}
