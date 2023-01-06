package jwt

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/home", handleHome)

	log.Panic(http.ListenAndServe(":9090", nil))
}

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		return
	}

}
