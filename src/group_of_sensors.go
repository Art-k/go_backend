package src

import (
	// "encoding/json"
	// "fmt"
	// "log"
	"fmt"
	"log"
	"net/http"
	// "strconv"
)

// SGroups list of groups
func SGroups(w http.ResponseWriter, r *http.Request) {

}

// SGroupCRUD crud for one group
func SGroupCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
	case "POST":
	case "PATCH":
	case "DELETE":
	case "GET":
	default:
		fmt.Fprintf(w, "Sorry, only OPTIONS,GET,POST,PATCH methods are supported. '"+r.Method+"' received")
		log.Println("/todo PATCH done\n\n")
	}
}
