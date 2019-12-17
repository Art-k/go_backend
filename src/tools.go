package src

import (
	"fmt"
	"net/http"
	"strings"

	guuid "github.com/satori/go.uuid"
)

/*

telegram bot key
BOT_API_KEY 					1041637917:AAHSWSqz-TfMUStJR88d3hrS8EnVtdoytF4
MY_CHANNEL_NAME					@father_fazenda
Text ????
Get request https://api.telegram.org/bot[BOT_API_KEY]/sendMessage?chat_id=[MY_CHANNEL_NAME]&text=[MY_MESSAGE_TEXT]

*/

const teleBotID = "1041637917:AAHSWSqz-TfMUStJR88d3hrS8EnVtdoytF4"
const teleBotChannel = "@father_fazenda"

func postTelegrammMessage(msg string) {

	var url string
	// fmt.Println(msg)
	url = "https://api.telegram.org/bot" + teleBotID + "/sendMessage?chat_id=" + teleBotChannel + "&parse_mode=HTML&text="

	msg = strings.Replace(msg, " ", "+", -1)
	msg = strings.Replace(msg, "'", "%27", -1)
	msg = strings.Replace(msg, "\n", "%0A", -1)

	url = url + msg
	fmt.Println("\n" + url + "\n")
	http.Get(url)

}

// GetHash we use it to get hasj=h for todo command
func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

// CheckIfTenOn get lates 5 values if see the difference
func CheckIfTenOn(incData IncomingDataStructure) {

	var currentState string

	var senseData []SenseDataTable
	Db.Where(&SenseDataTable{Mac: incData.Mac, Type: incData.Valuetype}).Limit(3).Order("created_at desc").Find(&senseData)
	for _, data := range senseData {
		fmt.Println(data)
	}

	if senseData[0].Value-senseData[len(senseData)-1].Value >= 2 {
		currentState = "ON"
	} else {
		Db.Where(&SenseDataTable{Mac: incData.Mac, Type: incData.Valuetype}).Limit(20).Order("created_at desc").Find(&senseData)
		for _, data := range senseData {
			fmt.Println(data)
		}
		if senseData[0].Value-senseData[len(senseData)-1].Value <= -2 {
			currentState = "OFF"
		}
	}

	if currentState != "" {
		fmt.Println(currentState)
		var ds DeviceState
		Db.Where("by_mac = ?", incData.Mac).Last(&ds)
		if ds.NewState != currentState {
			Db.Create(&DeviceState{ByMac: incData.Mac, NewState: currentState})
			postTelegrammMessage("Состояние изменено. текущее состояние : " + currentState)
		} else {
			fmt.Println("The same state")
		}
	}

}
