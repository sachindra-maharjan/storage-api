package events

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", status).Methods("GET")
	return router
}

func status(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Service is up and running!!")
}

