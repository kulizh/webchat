package main

import (
	"context"
	"log"
	"net/http"

	"webchat/internal/globals"
	"webchat/internal/routes"
)

func main() {
	backgroundContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	routes.Handle(backgroundContext)

	err := http.ListenAndServe(":"+globals.Config.ListenPort, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
