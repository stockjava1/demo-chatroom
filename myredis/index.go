package myredis

import (
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var (
	redisPool *redis.Pool
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	log *logger.CustZeroLogger
)

func init() {
	redisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Viper.GetString("redis.addr"))
		},
	}

	http.HandleFunc("/ws", handleWebsocket)
	log = logger.NewLoggerModule("redis")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Msgf("Fail to run http server %v", err)
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal().Msgf("Fail to upgrade websocket %v", err)
		return
	}

	go broadcastMessages(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error().Msgf("Fail to read message %v", err)
			return
		}

		err = saveMessageToRedis(string(message))
		if err != nil {
			log.Error().Msgf("Fail to save message to redis %v", err)
		}
	}
}

func broadcastMessages(conn *websocket.Conn) {
	pubsubConn := redisPool.Get()
	defer pubsubConn.Close()

	pubsubConn.Do("SUBSCRIBE", "chat")

	for {
		message, err := redis.String(pubsubConn.Receive())
		if err != nil {
			log.Error().Msgf("Fail to receive message from redis %v", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Error().Msgf("Fail to write message to redis %v", err)
			return
		}
	}
}

func saveMessageToRedis(message string) error {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("PUBLISH", "chat", message)
	if err != nil {
		return err
	}

	return nil
}
