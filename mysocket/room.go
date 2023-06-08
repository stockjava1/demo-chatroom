package mysocket

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var rooms map[string]*Room = make(map[string]*Room)
var roomsMutex sync.RWMutex

const (
	MAX_USERS = 20

	ERROR_PREFIX    = "Error: "
	ERROR_SEND      = ERROR_PREFIX + "You Need to Join/Create a Chat Room to send a message. \n"
	ERROR_CREATE    = ERROR_PREFIX + "Duplicate Name of Chat Room. Enter another name\n"
	ERROR_JOIN      = ERROR_PREFIX + "No Such ChatRoom Exits\n"
	ERROR_SWITCH    = ERROR_PREFIX + "No Such ChatRoom Exits\n A New Chat Room is creating..."
	ERROR_LEAVE     = ERROR_PREFIX + "You are out of the chatroom already.\n"
	ERROR_CLIENT    = ERROR_PREFIX + "Internal Server Error. Please Try Again!!!.\n"
	ERROR_NO_CLIENT = ERROR_PREFIX + "No such user exits. Internal Server Error. Plese disconnect the applicationa and start again!!!\n"

	ERROR_CONNECTION_CLOSE_FAIL = ERROR_PREFIX + "Connection Close Error. \n"

	NOTICE_ROOM_JOIN       = "\"%s\" joined\n"
	NOTICE_ROOM_LEAVE      = "\"%s\" left\n"
	NOTICE_ROOM_DELETE     = "Chat room is found inactive for seven days and deleted.\n"
	NOTICE_PERSONAL_CREATE = "Welcome to Chat Room \"%s\".\n"

	EXPIRY_TIME time.Duration = 7 * 24 * time.Hour
)

type Room struct {
	Name     string
	Clients  []*Client
	Messages []string
	Join     chan *Client
	Leave    chan *Client
	Incoming chan string
	Expire   chan bool
	Expiry   time.Time
}

func NewRoom(name string) *Room {
	room := &Room{
		Name:     name,
		Clients:  make([]*Client, 0),
		Messages: make([]string, 0),
		Join:     make(chan *Client),
		Leave:    make(chan *Client),
		Incoming: make(chan string),
		Expire:   make(chan bool),
		Expiry:   time.Now().Add(EXPIRY_TIME),
	}
	room.Listen()
	room.TryDelete()
	return room
}

func AddRoom(room *Room) error {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	otherRoom := rooms[room.Name]
	if otherRoom != nil {
		return errors.New(ERROR_CREATE)
	}
	rooms[room.Name] = room
	return nil
}

func RemoveRoom(name string) error {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room := rooms[name]
	if room == nil {
		return errors.New(ERROR_JOIN)
	}
	delete(rooms, name)
	return nil
}

func GetRoomNames() []string {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	keys := make([]string, 0, len(rooms))
	for k := range rooms {
		keys = append(keys, k)
	}
	return keys
}

func GetRoom(name string) (*Room, error) {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	room := rooms[name]
	if room == nil {
		return nil, errors.New(ERROR_JOIN)
	}
	return room, nil
}

func FindRoom(name string) (*Room, error) {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	room := rooms[name]
	if room == nil {
		return nil, errors.New(ERROR_SWITCH)
	}
	return room, nil
}

func (room *Room) AddClient(client *Client) {
	log.Println("AddClient")
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	room.Broadcast(fmt.Sprintf(NOTICE_ROOM_JOIN, client.Name))
	for _, message := range room.Messages {
		client.Outgoing <- message
	}
	room.Clients = append(room.Clients, client)
	client.Room = room
}

func (room *Room) RemoveClient(client *Client) {
	log.Println("RemoveConn")
	client.Mutex.RLock()
	room.Broadcast(fmt.Sprintf(NOTICE_ROOM_LEAVE, client.Name))
	client.Mutex.RUnlock()
	for i, otherClient := range room.Clients {
		if client == otherClient {
			room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
			break
		}
	}
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	client.Room = nil
}

func (room *Room) Broadcast(message string) {
	log.Println("Broadcast")
	room.Expiry = time.Now().Add(EXPIRY_TIME)
	log.Println(message)
	room.Messages = append(room.Messages, message)
	for _, client := range room.Clients {
		client.Outgoing <- message
	}
}

func (room *Room) Listen() {
	go func() {
		for {
			select {
			case message := <-room.Incoming:
				room.Broadcast(message)
			case client := <-room.Join:
				room.AddClient(client)
			case client := <-room.Leave:
				room.RemoveClient(client)
			case _ = <-room.Expire:
				room.TryDelete()
			}
		}
	}()
}

func (room *Room) TryDelete() {
	log.Println("TryDelete")
	if room.Expiry.After(time.Now()) {
		go func() {
			time.Sleep(room.Expiry.Sub(time.Now()))
			room.Expire <- true
		}()
	} else {
		room.Broadcast(NOTICE_ROOM_DELETE)
		for _, client := range room.Clients {
			client.Mutex.Lock()
			client.Room = nil
			client.Mutex.Unlock()
		}
		RemoveRoom(room.Name)
		//TODO: Clear out channels
	}
}
