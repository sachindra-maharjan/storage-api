// Package main is typicall used to define a single, executable command and its services

package main

import (
	"datasync/services/events"
	"log"
	"sync"
)

func main() {
	log.Println("Starging all services")

	// start the servers
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go startServer(wg, events.StartServer)
	wg.Wait()
	log.Println("All services stopped")
}

// startServer is a function in the main package that supports the command.
func startServer(wg *sync.WaitGroup, startFunc func() error) {
	err := startFunc()
	wg.Done()
	if err != nil {
		log.Fatal(err)
	}
}