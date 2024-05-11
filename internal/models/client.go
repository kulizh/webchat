package models

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type ManagerInterface interface {
	AddClient(client *Client)
	RemoveClient(client *Client)
	GetList() ClientList
}

type Client struct {
	Connection *websocket.Conn
	egress     chan Event
	chatroom   string
	handlers   map[string]EventHandler
	manager    ManagerInterface
	sync.RWMutex
}

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

var (
	ErrEventNotSupported = errors.New("Event type not supported")
)

func NewClient(conn *websocket.Conn, m ManagerInterface) *Client {
	client := &Client{
		Connection: conn,
		egress:     make(chan Event),
		handlers:   make(map[string]EventHandler),
		manager:    m,
	}

	client.setupEventHandlers()
	client.add()

	return client
}

func (client *Client) ReadMessages() {
	defer func() {
		client.remove()
	}()

	client.Connection.SetReadLimit(512)

	if err := client.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	client.Connection.SetPongHandler(client.pongHandler)

	for {
		_, payload, err := client.Connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break
		}

		if err := client.RouteEvent(request); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

func (client *Client) remove() {
	client.Lock()
	defer client.Unlock()

	client.manager.RemoveClient(client)
}

func (client *Client) add() {

	client.Lock()
	defer client.Unlock()

	client.manager.AddClient(client)

	log.Println("New client")
}

func (client *Client) pongHandler(pongMsg string) error {
	return client.Connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *Client) WriteMessages() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		client.remove()
	}()

	for {
		select {
		case message, ok := <-client.egress:

			if !ok {
				if err := client.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ", err)
				}

				return
			}

			data, err := json.Marshal(message)

			if err != nil {
				log.Println(err)

				return
			}

			if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}

		case <-ticker.C:
			log.Println("ping")

			if err := client.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return
			}
		}
	}
}

func (client *Client) RouteEvent(event Event) error {
	if handler, ok := client.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (client *Client) setupEventHandlers() {
	client.handlers[EventSendMessage] = SendMessageHandler
	client.handlers[EventChangeRoom] = ChatRoomHandler
}

func (client *Client) GetList() ClientList {
	return client.manager.GetList()
}
