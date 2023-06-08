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

type Client struct {
	Token       string
	Name        string
	UID         string
	ID          string
	Room        *Room
	Outgoing    chan string
	Mutex       sync.RWMutex
	conn        *neffos.Conn
	Logined     bool
	log         *logger.CustZeroLogger
	ConnectTime time.Time
}

func NewClient(uid string, conn *neffos.Conn, log *logger.CustZeroLogger) *Client {
	return &Client{
		Token:       "",
		Name:        "",
		Room:        nil,
		UID:         uid,
		conn:        conn,
		Outgoing:    make(chan string),
		Logined:     false,
		ConnectTime: time.Now(),
		log:         log,
	}
}

func (c *Client) Close() error {
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
		c.log.Error().Msgf("Send user %s , conn %s fail", c.UID, c.conn.ID())
	}
}

func (c *Client) HandleEvent(event string, data interface{}) {
	c.log.Info().Msgf("handle event %s, data %v", event, data)

	switch event {
	case "userInfo":
		c.getUserInfo()
		break
	default:
		//
	}
}
func (c *Client) getUserInfo() {
	userInfo, err := json.Marshal(&response.UserInfo{c.UID, c.Name, c.ID})
	if err == nil {
		c.Send(string(userInfo))
	}
}
