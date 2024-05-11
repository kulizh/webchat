package service

import (
	"context"
	"log"
	"net/http"

	"webchat/internal/models"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func checkOrigin(r *http.Request) bool {
	return true

	//origin := r.Header.Get("Origin")

	// switch origin {
	// case "http://localhost:8080":
	// 	return true
	// default:
	// 	return false
	//}
}

type Manager struct {
	clients models.ClientList
}

func NewManager(backgroundContext context.Context) *Manager {
	return &Manager{
		clients: make(models.ClientList),
	}
}

func (manager *Manager) ServeSockets(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create New Client
	client := models.NewClient(conn, manager)

	go client.ReadMessages()
	go client.WriteMessages()
}

func (manager *Manager) AddClient(client *models.Client) {
	manager.clients[client] = true
}

func (manager *Manager) RemoveClient(client *models.Client) {
	if _, ok := manager.clients[client]; ok {
		client.Connection.Close()
		delete(manager.clients, client)
	}
}

func (manager *Manager) GetList() models.ClientList {
	return manager.clients
}
