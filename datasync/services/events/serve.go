package events

import (
	"datasync/services/events/internals/ports"
	"log"
	"net/http"
	"os"
)

var port string = "80"

func StartServer() error{
	router := Routes()
	log.Printf("Running on port: %s", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = string(ports.EventService)
		log.Printf("Defaulting to port %s", port)
	}	
}