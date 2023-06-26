package mysocket

import (
	"encoding/json"
	"errors"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/JabinGP/demo-chatroom/mysocket/response"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"sync"
	"time"
)

var Clients map[string]*Client = make(map[string]*Client)
var UIDs map[string]map[string]string = make(map[string]map[string]string) //map的value存储uid，用于区分用户 及对应的ClientId

var ClientsMutex sync.RWMutex

type Client struct {
	Token       string
	Name        string
	UID         string
	ID          string
	CID         string
	Rooms       map[string]string
	Outgoing    chan string
	Mutex       sync.RWMutex
	conn        *neffos.Conn
	Logined     bool
	log         *logger.CustZeroLogger
	ConnectTime time.Time
	Publish     func(string, string)
}

func NewClient(id string, cid string, conn *neffos.Conn, log *logger.CustZeroLogger, publish func(string, string)) *Client {
	return &Client{
		Token:       "",
		Name:        "",
		Rooms:       make(map[string]string),
		ID:          id,
		CID:         cid,
		conn:        conn,
		Outgoing:    make(chan string),
		Logined:     false,
		ConnectTime: time.Now(),
		log:         log,
		Publish:     publish,
	}
}

func (c *Client) Close() error {
	ClientsMutex.RLock()
	defer ClientsMutex.RUnlock()
	for _, roomId := range c.Rooms {
		if r, ok := Rooms[roomId]; ok {
			if _, ok := r.ClientIDs[c.ID]; ok {
				delete(r.ClientIDs, c.ID)
			}
		}
	}

	if _, ok := UIDs[c.UID]; ok {
		delete(UIDs[c.UID], c.ID)
		if len(UIDs[c.UID]) == 0 {
			delete(UIDs, c.UID)
		}
	}

	if _, ok := Clients[c.ID]; ok {
		delete(Clients, c.ID)
	}

	error := c.conn.Socket().NetConn().Close()
	c.log.Info().Msgf("Try close socket. id = %s, uid = %s", c.conn.ID(), c.UID)
	if error != nil {
		c.log.Error().Msgf("Close socket error. id = %s, uid = %s", c.conn.ID(), c.UID)

		return errors.New(ERROR_CONNECTION_CLOSE_FAIL)
	}

	return nil
}

func (c *Client) Send(msg string) {
	ok := c.conn.Write(websocket.Message{
		Body:     []byte(msg), // []byte{byte(1)},
		IsNative: true,
	})
	if !ok {
		c.log.Error().Msgf("Send user %s , cid %s, client id %s, conn %s fail", c.UID, c.CID, c.ID, c.conn.ID())
	}
}

func (c *Client) Login(userId string, userName string) error {
	c.UID = userId
	c.Name = userName

	//c.Name = tokenString["userName"].(string)
	//c.Token = token.(string)
	c.Logined = true

	if _, ok := UIDs[c.UID]; !ok {
		UIDs[c.UID] = make(map[string]string)
	}
	if _, ok := UIDs[c.UID][c.ID]; !ok {
		UIDs[c.UID][c.ID] = c.ID
	}

	return nil
}

func (c *Client) JoinRoom(roomId string) {
	var r *Room
	if r, ok := Rooms[roomId]; !ok {
		r = NewRoom(roomId, "Room-"+roomId)
		Rooms[roomId] = r
	}
	r = Rooms[roomId]
	//if room, ok := Rooms[roomId]; ok {
	//	Rooms[roomId] = NewRoom(roomId, "Room-"+roomId)
	//}

	if _, ok := r.ClientIDs[c.ID]; !ok {
		r.ClientIDs[c.ID] = c.ID
	}
}

func (c *Client) HandleEvent(event string, data interface{}) error {
	c.log.Info().Msgf("handle event %s, data %v", event, data)

	switch event {
	case "login":
		if d, ok := data.(map[string]interface{}); ok {
			name := d["name"].(string)
			id := d["id"].(string)
			c.Login(id, name)
		} else {
			return errors.New("login fail")
		}
		break
	case "userInfo":
		c.getUserInfo()
		break
	case "userMsg":
		if d, ok := data.(map[string]interface{}); ok {
			uid := d["uid"].(string)
			msg := d["msg"].(string)
			c.Publish(uid, msg)
		}
		break
	case "joinRoom":
		if d, ok := data.(map[string]interface{}); ok {
			roomId := d["roomId"].(string)
			c.JoinRoom(roomId)
		}
		break
	case "roomMsg":
		if d, ok := data.(map[string]interface{}); ok {
			roomId := d["roomId"].(string)
			msg := d["msg"].(string)
			c.log.Info().Msgf("room ID: %s, msg: %s, rooms: %v", roomId, msg, Rooms)
			if r, exist := Rooms[roomId]; exist {
				c.log.Info().Msgf("room ID: %s, msg: %s, clients: %v", roomId, msg, Rooms[roomId])
				for _, clientId := range r.ClientIDs {
					if client, ok := Clients[clientId]; ok {
						c.log.Info().Msgf("Send msg to client id:%s, cid:%s, userid:%s, username:%s, room ID: %s, msg: %s", client.ID, client.CID, client.UID, client.Name, roomId, msg)
						client.Send(msg)
					}
				}
			}
		}
		break
	default:
		c.Send("invalid request")
	}

	return nil
}

func (c *Client) getUserInfo() {
	userInfo, err := json.Marshal(&response.UserInfo{c.UID, c.Name, c.ID})
	if err == nil {
		c.Send(string(userInfo))
	}
}
