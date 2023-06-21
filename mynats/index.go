package mynats

import (
	"fmt"
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/nats-io/nats.go"
	"time"
)

var log *logger.CustZeroLogger

var NatsConn *nats.Conn

var timeout = 5 * time.Second

func init() {

}

var subs map[string]*nats.Subscription
var natsConn *nats.Conn

func init() {
	log = logger.NewLoggerModule("nats")
	// 连接NATS服务器
	natsConn, err := nats.Connect(fmt.Sprintf("nats://%s", config.Viper.GetString("nats.addr")))

	//defer natsConn.Close()
	NatsConn = natsConn

	if err != nil {
		log.Fatal().Msgf("Fail to connect nats %v", err)
		panic("Fail to connect nats %v")
	}
	// 创建一个信号通道，以便在收到操作系统信号时优雅地关闭连接
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	subs = make(map[string]*nats.Subscription)
	healthCheck()

	// 等待操作系统信号
	// <-sigChan

	//log.Info().Msgf("Disconnect nats")

}

func Unsubscribe() {
	for uid, sub := range subs {
		err := sub.Unsubscribe()
		if err != nil {
			log.Error().Msgf("Error unsubscribing user %s: %v", uid, err)
		} else {
			log.Info().Msgf("unsubscribed user %s: %v", uid)
		}
		delete(subs, uid)
	}
	natsConn.Close()
	log.Info().Msgf("Disconnect nats")
}

func Subscribe(uid string) {
	// defer nc.Close()

	// 订阅NATS消息
	sub, err := NatsConn.Subscribe(uid, func(msg *nats.Msg) {
		log.Info().Msgf("收到主题 %s 的消息：%s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Error().Msgf("订阅主题 %s 失败：%v\n", uid, err)
	} else {
		log.Info().Msgf("成功订阅主题 %s\n", uid)
		subs[uid] = sub
	}
}

func healthCheck() {
	id := "healthCheck"
	Subscribe(id)
	Publish(id, "i'm healthy!")
}
func Publish(uid string, msg string) {
	message := []byte(msg)
	err := NatsConn.Publish(uid, message)
	if err != nil {
		log.Error().Msgf("Fail to send nats uid %s, msg %s, err %v", uid, msg, err)
	} else {
		log.Info().Msgf("==> send nats send nats uid %s, msg %s, err %v", uid, msg)
	}
}
