package routes

import (
	"context"
	"fmt"
	"net/http"

	"webchat/internal/globals"
	"webchat/internal/service"
)

func Handle(backgroundContext context.Context) {
	manager := service.NewManager(backgroundContext)

	//http.Handle("/", http.FileServer(http.Dir("../../webchat")))

	http.HandleFunc(globals.Config.Routes.Ws, manager.ServeSockets)

	http.HandleFunc(globals.Config.Routes.Counter, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, len(manager.GetList()))
	})
}
