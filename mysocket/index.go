package mysocket

import (
	"encoding/json"
	"errors"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"sync"
	"time"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate
var clientsMutex sync.RWMutex

func init() {
	validate = validator.New()
}

type MyWebSocket struct {
	Ws    *neffos.Server
	Conns map[*neffos.Conn]*Client      //map的value存储uid，用于区分用户
	IDs   map[string]map[string]*Client //map的value存储uid，用于区分用户 及对应的Client
	log   *logger.CustZeroLogger
	//NatsConn *nats.Conn //map的value存储cilent id，用于区分不同client 对应的connect
}

func NewSocket() *MyWebSocket {
	log := logger.NewLoggerModule("socket")
	//conn, err := mynats.Conn()
	//if err != nil {
	//	log.Fatal().Msgf("Fail to connect nats %v", err)
	//	panic("Fail to connect nats")
	//}

	mySocket := MyWebSocket{
		Conns: make(map[*neffos.Conn]*Client),
		IDs:   make(map[string]map[string]*Client),
		log:   log,
		//NatsConn: conn,
	}

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			mySocket.log.Info().Msgf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
			client := mySocket.Conns[nsConn.Conn]
			body := msg.Body

			var eventData interface{}
			err := json.Unmarshal(body, &eventData)
			if err != nil {
				mySocket.log.Error().Msgf("Error %v", err)
			} else {
				mySocket.log.Info().Msgf("Event %v", eventData)
				// ping pong
				if eventData == "2" {
					mySocket.log.Info().Msgf("Pong %v", eventData)
					return nil
				}

				// 使用类型断言来判断 data 的类型是否为 map[string]interface{}
				if m, ok := eventData.(map[string]interface{}); ok {
					// 如果是 map[string]interface{} 类型，说明 JSON 数据是一个对象
					// 可以使用 m["event"] 和 m["data"] 来获取 event 和 data 的值
					eventName := m["event"]
					data := m["data"]

					if eventName == "login" {
						// 假设 m["data"] 的值不为 nil，并且是 map[string]interface{} 类型
						if d, ok := data.(map[string]interface{}); ok {
							// 如果是 map[string]interface{} 类型，说明 data 是一个对象
							// 可以使用 d["name"] 和 d["age"] 来获取 name 和 age 的值

							/*
								token := d["token"]
								mySocket.log.Info().Msgf("Token %s", token)
								tokenString, err := tool.ParseToken(token.(string))
								if err != nil {
									mySocket.log.Error().Msgf("Parse Token error %v", err)
									mySocket.RemoveConn(client.conn)
									return nil
								}

								mySocket.log.Info().Msgf("tokenString %v", tokenString)
							*/
							userName := d["userName"].(string)
							userID := d["userId"].(string)

							res := "login success"
							client.Name = userName
							client.ID = userID
							//client.Name = tokenString["userName"].(string)
							//client.Token = token.(string)
							client.Logined = true
							//client.ID = tokenString["userId"].(string)
							mySocket.AddClient(client)
							client.SubscribeMsg()
							client.Send(res)
							mySocket.log.Info().Msgf("uid [%s] Login!", mySocket.Conns[nsConn.Conn])
						} else {
							// 如果不是 map[string]interface{} 类型，说明 data 不是一个对象
							// 可以打印 m["data"] 的类型和值
							mySocket.log.Error().Msgf("fail to login. close connection")
							mySocket.RemoveConn(client.conn)
							return nil
						}
					} else {

						client := mySocket.Conns[nsConn.Conn]
						if !client.Logined {
							client.Send("Please login before access")
							mySocket.log.Info().Msgf("Please login before access")
							return nil
						}

						if d, ok := data.(map[string]interface{}); ok {
							client.HandleEvent(eventName.(string), d)
						} else {
							mySocket.log.Info().Msgf("unhandled data %s %v", eventName, d)
						}

					}

				} else {
					// 如果不是 map[string]interface{} 类型，说明 JSON 数据不是一个对象
					// 可以打印 data 的类型和值
					mySocket.log.Info().Msgf("Received data %T %v", eventData, eventData)
				}

			}

			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		ctx := websocket.GetContext(c)
		uid := ctx.URLParam("uid")
		mySocket.log.Info().Msgf("[%s] user %s Connected to server!", c.ID(), uid)

		//mySocket.SetUID(c, uid)
		//client := NewClient(token.(string), tokenString["userName"].(string), tokenString["userId"].(string), nsConn.Conn)
		client := NewClient(uid, c, mySocket.log)
		mySocket.AddClient(client)

		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		mySocket.RemoveConn(c)
		mySocket.log.Info().Msgf("[%s] Disconnected from server", c.ID())
	}

	ws.OnUpgradeError = func(err error) {
		mySocket.log.Info().Msgf("Upgrade Error: %v", err)
	}

	mySocket.Ws = ws

	return &mySocket
}

// SetUID 设置用户信息
func (m *MyWebSocket) AddClient(client *Client) error {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	otherClients := m.Conns[client.conn]
	if otherClients != nil {
		return errors.New(ERROR_CLIENT)
	}
	m.log.Info().Msgf("Set connect [%s] uid: [%s]", client.conn.ID(), client.UID)
	m.Conns[client.conn] = client

	if client.Logined {
		_, ok := m.IDs[client.UID]
		if !ok { // key不存在
			m.IDs[client.UID] = make(map[string]*Client)
		}
		m.IDs[client.UID][client.ID] = client
	}
	return nil
}

// SetUID 设置用户信息
func (m *MyWebSocket) GetClient(c *neffos.Conn) (*Client, error) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	client := m.Conns[c]
	if client == nil {
		return nil, errors.New(ERROR_NO_CLIENT)
	}
	return client, nil
}

// SetUID 设置用户信息
func (m *MyWebSocket) GetClientByUID(uid string) (map[string]*Client, error) {
	//clientsMutex.Lock()
	//defer clientsMutex.Unlock()

	if _, ok := m.IDs[uid]; ok {
		return m.IDs[uid], nil
	} else {
		return nil, errors.New("NO_USER_CLIENT")
	}

}

// DelConn 移除连接
func (m *MyWebSocket) RemoveConn(c *websocket.Conn) error {
	m.log.Info().Msgf("delete connect [%s] uid: [%s]", c.ID(), m.Conns[c])

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	client := m.Conns[c]

	if client == nil {
		return errors.New(ERROR_NO_CLIENT)
	}

	delete(m.IDs, client.UID)
	delete(m.Conns, c)

	return nil
}

func (m *MyWebSocket) Ping() {
	pingTicker := time.NewTicker(10 * time.Second)

	go func() {
		defer pingTicker.Stop()

		for {
			select {
			case <-pingTicker.C:
				//m.log.Info().Msgf("ticker conns %d", len(m.Conns))
				for _, client := range m.Conns {
					//m.log.Info().Msgf("ping %s", client.conn.ID())
					client.Send("1")
				}
			default:
			}
		}
	}()
}
