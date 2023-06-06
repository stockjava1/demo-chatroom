package mysocket

import (
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

type MyWebSocket struct {
	Ws    *neffos.Server
	Conns map[*neffos.Conn]string //map的value存储uid，用于区分用户
	log   *logger.CustZeroLogger
}

func NewSocket() *MyWebSocket {
	mySocket := MyWebSocket{
		Conns: make(map[*neffos.Conn]string),
		log:   logger.NewLoggerModule("socket"),
	}

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			mySocket.log.Info().Msgf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			ping := string(msg.Body)
			pong := strings.Replace(ping, "？", "！", len(ping))
			pong = strings.Replace(pong, "么", "", len(pong))

			mg := websocket.Message{
				Body:     []byte(pong),
				IsNative: true,
			}

			nsConn.Conn.Write(mg)

			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		mySocket.log.Info().Msgf("[%s] Connected to server!", c.ID())
		ctx := websocket.GetContext(c)
		uid := ctx.URLParam("uid")
		mySocket.log.Info().Msgf("uid [%s] Connected to server!", uid)
		mySocket.SetUID(c, uid)
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		mySocket.DelConn(c)
		mySocket.log.Info().Msgf("[%s] Disconnected from server", c.ID())
	}

	ws.OnUpgradeError = func(err error) {
		mySocket.log.Info().Msgf("Upgrade Error: %v", err)
	}

	mySocket.Ws = ws

	return &mySocket
}

// SetUID 设置用户信息
func (m *MyWebSocket) SetUID(c *neffos.Conn, uid string) {
	m.log.Info().Msgf("Set connect [%s] uid: [%s]", c.ID(), uid)
	m.Conns[c] = uid
}

// DelConn 移除连接
func (m *MyWebSocket) DelConn(c *neffos.Conn) {
	m.log.Info().Msgf("delete connect [%s] uid: [%s]", c.ID(), m.Conns[c])
	delete(m.Conns, c)
}
