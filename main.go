package main

import (
	"fmt"
	"github.com/JabinGP/demo-chatroom/config"
	_ "github.com/JabinGP/demo-chatroom/docs"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/JabinGP/demo-chatroom/middleware"
	"github.com/JabinGP/demo-chatroom/mynats"
	"github.com/JabinGP/demo-chatroom/mysocket"
	"github.com/JabinGP/demo-chatroom/route"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// proxyResponseWriter is a wrapper around http.ResponseWriter that captures the status code and size
type proxyResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader captures the status code and calls the underlying WriteHeader method
func (prw *proxyResponseWriter) WriteHeader(code int) {
	prw.statusCode = code
	prw.ResponseWriter.WriteHeader(code)
}

// Write captures the size of the data and calls the underlying Write method
func (prw *proxyResponseWriter) Write(data []byte) (int, error) {
	size, err := prw.ResponseWriter.Write(data)
	prw.size += size
	return size, err
}

// irisZerologMiddleware 是一个使用 zerolog 的 iris 中间件
func irisZerologMiddleware(ctx iris.Context) {
	// 获取请求的方法和路径
	method := ctx.Method()
	path := ctx.Path()

	// 创建一个子日志器，添加请求的字段
	zlogger := logger.NewLoggerModule("http")
	subLogger := zlogger.With().Str("component", "web").
		Str("method", method).
		Str("path", path).
		Logger()
	zlogger.SetLogger(&subLogger)

	// 记录请求开始的日志
	zlogger.Info().Msg("request started")

	// 创建一个代理响应写入器，捕获状态码和大小
	prw := &proxyResponseWriter{ctx.ResponseWriter(), http.StatusOK, 0}

	// 测量处理请求所需的时间
	start := time.Now()
	ctx.Next()
	elapsed := time.Since(start)

	// 记录请求结束的日志
	//output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	//zlog2 := zlog.Output(output).With().Timestamp().Logger()
	//zlog2.Info().
	//	Int("status", prw.statusCode).
	//	Int("size", prw.size).
	//	Dur("elapsed", elapsed).
	//	Msg("request completed")

	zlogger.Info().Msgf("status %d, size %d, elapsed %v, request completed", prw.statusCode, prw.size, elapsed)
}

type MyServer struct {
	MySocket *mysocket.MyWebSocket
	MyNats   *mynats.MyNats
}

// swagger middleware for Iris
// swagger embed files

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 初始化 Logger
	// logger.InitLog()
	// 在 main 包中使用 Logger

	log := logger.NewLoggerModule("app")
	log.Info().Msg("start application")

	//logLevel := config.Viper.GetString("server.logger.level")
	//logger.SetLogLevel(logLevel)

	mySocket := mysocket.NewSocket()
	mySocket.Ping()

	mySocket.NatsCheck()

	app := iris.New()

	app.Get("/ws", websocket.Handler(mySocket.Ws))

	//app.Logger().SetLevel(config.Viper.GetString("server.logger.level"))
	// Add logger to log the requests to the terminal
	//app.Use(logger.New())
	// 使用 zerolog 中间件
	app.Use(irisZerologMiddleware)

	// Add recover to recover from any http-relative panics
	app.Use(recover.New())
	// Globally allow options method to enable CORS
	app.AllowMethods(iris.MethodOptions)
	// Add global CORS handler
	app.Use(middleware.CORS)

	// Router
	route.Route(app)

	// Listen in 8888 port
	app.Run(iris.Addr(config.Viper.GetString("server.addr")), iris.WithoutServerError(iris.ErrServerClosed))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	// 等待操作系统信号
	//<-sigChan

	// 这里是程序的主循环
	for {
		select {
		case <-sigChan:
			// 接收到 os.Interrupt 信号，执行程序退出操作
			fmt.Println("Received interrupt signal, exiting...")
			mySocket.UnsubscribeAll()
			os.Exit(0)
		default:
			// 程序的主要逻辑
			//fmt.Println("Running...")
		}
	}
	//
}
