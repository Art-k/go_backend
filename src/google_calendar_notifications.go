package src

import (
	"fmt"
	"log"
	"net/http"
)

type apiChannelStructure struct {
	kind        string
	id          string
	resourceId  string
	resourceUri string
	token       string
	expiration  int32
}

func Notifications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		OptionsAnswer(w, r)
		return
	case "POST":

		log.Println("Notification Received")

		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, "")
		fmt.Println(n)

		log.Println("Notification Done")

	case "PATCH":
	case "GET":
	case "DELETE":
	default:
	}
}
