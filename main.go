package main

import (
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/JabinGP/demo-chatroom/middleware"
	"github.com/JabinGP/demo-chatroom/route"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
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
	log := logger.NewLogger()
	log.SetLogLevel(config.Viper.GetString("loglevel.http"))
	log.SetModule("http")
	subLogger := log.GetLogger().With().Str("component", "web").
		Str("method", method).
		Str("path", path).
		Logger()
	log.SetLogger(&subLogger)
	// 记录请求开始的日志
	log.Info("request started")

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

	log.Info("status %d, size %d, elapsed %v, request completed", prw.statusCode, prw.size, elapsed)
}

func main() {
	// 初始化 Logger
	//logger.InitLog()
	// 在 main 包中使用 Logger

	//logLevel := config.Viper.GetString("server.logger.level")
	//logger.SetLogLevel(logLevel)
	log := logger.NewLogger()
	log.SetLogLevel(config.Viper.GetString("loglevel.app"))
	log.SetModule("app")
	log.Info("start application")
	app := iris.New()
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
}
