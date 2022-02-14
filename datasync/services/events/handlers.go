package events

import (
	"context"
	"datasync/services/events/data"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", status).Methods("GET")
	router.HandleFunc("/standings/{leagueId}", standings).Methods("GET")
	return router
}

func status(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Service is up and running!!")
}

func standings(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leagueId := vars["leagueId"]
		client := data.NewClient(nil, []string{"U4y3LniAIdmsh1SryySGibO7k8ELp1syFPvjsnpHOQNWAvpJAk"});
		id, err := strconv.Atoi(leagueId);
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		result, _, err := client.StandingService.GetLeagueStandings(context.Background(), id);
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(result)
}