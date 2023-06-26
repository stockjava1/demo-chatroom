package mysocket

import (
	"encoding/json"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/JabinGP/demo-chatroom/mynats"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"time"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type MyWebSocket struct {
	Ws     *neffos.Server
	Conns  map[*neffos.Conn]string //map的value存储uid，用于区分用户
	log    *logger.CustZeroLogger
	MyNats *mynats.MyNats
}

func NewSocket() *MyWebSocket {
	log := logger.NewLoggerModule("socket")
	//conn, err := mynats.Conn()
	//if err != nil {
	//	log.Fatal().Msgf("Fail to connect nats %v", err)
	//	panic("Fail to connect nats")
	//}

	myNats, err := mynats.NewMyNats()
	if err != nil {
		panic("Fail to init nats %v")
	}
	mySocket := MyWebSocket{
		Conns:  make(map[*neffos.Conn]string),
		log:    log,
		MyNats: myNats,
		//NatsConn: conn,
	}
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			mySocket.log.Info().Msgf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
			//clientId := mySocket.Conns[nsConn.Conn]
			client := Clients[nsConn.Conn.ID()]
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
							err := client.HandleEvent("login", d)
							if err != nil {

								mySocket.log.Error().Msgf(" [%s] Login fail!", data)
								res := "login fail"
								//client.ID = tokenString["userId"].(string)
								//mySocket.AddClient(client)
								//mySocket.Subscribe(client.UID)
								//client.Login()
								client.Emit(res)

							} else {
								mySocket.log.Info().Msgf("uid [%s] Login!", mySocket.Conns[nsConn.Conn])
								res := "login success"
								//client.ID = tokenString["userId"].(string)
								//mySocket.AddClient(client)
								mySocket.SubscribeUser(client.UID)
								//client.Login()
								client.Emit(res)
								mySocket.RemoveConn(nsConn.Conn)
							}

						} else {
							// 如果不是 map[string]interface{} 类型，说明 data 不是一个对象
							// 可以打印 m["data"] 的类型和值
							mySocket.log.Error().Msgf("fail to login. close connection")
							mySocket.RemoveConn(nsConn.Conn)
							return nil
						}
					} else {

						//client := mySocket.Conns[nsConn.Conn]
						if !client.Logined {
							client.Emit("Please login before access")
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
		cid := ctx.URLParam("cid")
		mySocket.log.Info().Msgf("[%s] connect %s Connected to server!", c.ID(), cid)

		//mySocket.SetUID(c, uid)
		//client := NewClient(token.(string), tokenString["userName"].(string), tokenString["userId"].(string), nsConn.Conn)
		client := NewClient(c.ID(), cid, c, mySocket.log, mySocket.Publish)
		mySocket.Conns[c] = client.ID
		//mySocket.AddClient(client)
		userId := ctx.URLParam("userId")
		userName := ctx.URLParam("userName")
		roomId := ctx.URLParam("roomId")
		err := client.Login(userId, userName)
		if err == nil {
			res := "login success"
			mySocket.SubscribeUser(client.UID)
			//client.Login()
			client.Emit(res)

			client.JoinRoom(roomId)
			mySocket.log.Info().Msgf("[%s] user %s Connected to server!", c.ID(), cid)

			Clients[client.ID] = client
		} else {
			res := "login fail"
			client.Emit(res)
			c.Close()
		}
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

/*
// SetUID 设置用户信息
func (m *MyWebSocket) AddClient(client *Client) error {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	otherClients := m.Conns[client.conn]
	if otherClients != nil {
		return errors.New(ERROR_CLIENT)
	}
	m.log.Info().Msgf("Set connect [%s] uid: [%s]", client.conn.ID(), client.UID)
	m.Conns[client.conn] = client

	if client.Logined {
		_, ok := m.UIDs[client.UID]
		if !ok { // key不存在
			m.UIDs[client.UID] = make(map[string]*Client)
		}
		m.UIDs[client.UID][client.ID] = client
	}
	return nil
}
*/
// DelConn 移除连接
func (m *MyWebSocket) RemoveConn(c *websocket.Conn) error {
	m.log.Info().Msgf("delete connect [%s] uid: [%s]", c.ID(), m.Conns[c])

	ClientsMutex.RLock()
	defer ClientsMutex.RUnlock()

	if clientId, ok := m.Conns[c]; ok {
		if c, ok := Clients[clientId]; ok {
			delete(Clients, clientId)
			c.Close()
		}
	}
	//if _, ok := m.Conns[c]; ok {
	delete(m.Conns, c)
	//}

	return nil
}

func (m *MyWebSocket) NatsCheck() {
	id := "healthCheck"
	m.SubscribeUser(id)
	m.Publish(id, "i'm healthy!")
}

func (m *MyWebSocket) Publish(uid string, msg string) {
	message := []byte(msg)
	err := m.MyNats.NatsConn.Publish(uid, message)
	if err != nil {
		m.MyNats.Log.Error().Msgf("Fail to send nats uid %s, msg %s, err %v", uid, msg, err)
	} else {
		m.MyNats.Log.Info().Msgf("==> send nats send nats uid %s, msg %s, err %v", uid, msg)
	}
}

func (m *MyWebSocket) SubscribeUser(uid string) {
	// defer nc.Close()

	_, ok := m.MyNats.UserSubs[uid]
	if !ok {
		// 订阅NATS消息
		sub, err := m.MyNats.NatsConn.Subscribe(uid, func(msg *nats.Msg) {
			msgData := string(msg.Data)
			log.Info().Msgf("收到主题 %s 的消息：%s\n", msg.Subject, string(msg.Data))

			uClient, ok := UIDs[uid]
			if ok { // key不存在
				//m.UIDs[client.UID] = make(map[string]*Client)
				//map[string]*mysocket.Client
				for _, clientId := range uClient {
					if client, ok := Clients[clientId]; ok {
						client.Emit(msgData)
					}
				}
			}

		})
		if err != nil {
			log.Error().Msgf("订阅主题 %s 失败：%v\n", uid, err)
		} else {
			log.Info().Msgf("成功订阅主题 %s\n", uid)
			m.MyNats.UserSubs[uid] = sub
		}
	}

}

func (m *MyWebSocket) UnsubscribeAll() {
	m.MyNats.Unsubscribe()
}

func (m *MyWebSocket) Ping() {
	pingTicker := time.NewTicker(10 * time.Second)

	go func() {
		defer pingTicker.Stop()

		for {
			select {
			case <-pingTicker.C:
				//m.log.Info().Msgf("ticker conns %d", len(m.Conns))
				for _, clientId := range m.Conns {
					if c, ok := Clients[clientId]; ok {
						//m.log.Info().Msgf("ping %s", client.conn.ID())
						c.Emit("1")
					}
				}
			default:
			}
		}
	}()
}
