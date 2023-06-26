package mynats

import (
	"fmt"
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"time"
)

var timeout = 5 * time.Second

type MyNats struct {
	Subs     map[string]*nats.Subscription
	NatsConn *nats.Conn
	//UIDs      map[string]map[string]*mysocket.Client
	Log *logger.CustZeroLogger
}

func NewMyNats() (*MyNats, error) {
	log := logger.NewLoggerModule("nats")
	// 连接NATS服务器
	natsConn, err := nats.Connect(fmt.Sprintf("nats://%s", config.Viper.GetString("nats.addr")))

	if err != nil {
		log.Fatal().Msgf("Fail to connect nats %v", err)
		panic("Fail to connect nats %v")
	}
	// 创建一个信号通道，以便在收到操作系统信号时优雅地关闭连接
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	subs := make(map[string]*nats.Subscription)
	//UIDs := make(map[string]map[string]*mysocket.Client)
	myNats := &MyNats{
		Subs:     subs,
		NatsConn: natsConn,
		//UIDs:      UIDs,
		Log: log,
	}

	return myNats, nil

	// 等待操作系统信号
	// <-sigChan

	//log.Info().Msgf("Disconnect nats")

}

func (n *MyNats) Unsubscribe() {
	for uid, sub := range n.Subs {
		err := sub.Unsubscribe()
		if err != nil {
			n.Log.Error().Msgf("Error unsubscribing user %s: %v", uid, err)
		} else {
			n.Log.Info().Msgf("unsubscribed user %s: %v", uid)
		}
		delete(n.Subs, uid)
	}
	n.NatsConn.Close()
	log.Info().Msgf("Disconnect nats")
}
