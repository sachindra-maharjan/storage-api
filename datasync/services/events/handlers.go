package events

import (
	"context"
	"datasync/services/events/data/source"
	"datasync/services/events/data/store"
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
		client := source.NewClient(nil, []string{"U4y3LniAIdmsh1SryySGibO7k8ELp1syFPvjsnpHOQNWAvpJAk"});
		id, err := strconv.Atoi(leagueId);
		if err != nil {
			writeError(w, err)
			return
		}
		result, err := client.StandingService.GetLeagueStandings(context.Background(), id);
		if err != nil {
			writeError(w, err)
			return
		}

		fsClient, err := store.NewClient(context.Background(), "clouddeveloper-299318")
		if err != nil {
			writeError(w, err)
			return
		}

		err = fsClient.StandingsService.Set(context.Background(), "premierleague", result)
		

		json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(err)
}