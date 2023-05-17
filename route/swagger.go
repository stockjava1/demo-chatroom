package route

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12/core/router"
)

func routeSwagger(party router.Party) {
	/*
		// swagger配置方法一：其他文档
		config := swagger.Config{
			// 指向swagger init生成文档的路径
			URL:         "http://www.xxx.com/swagger/doc.json",
			DeepLinking: true,
		}
		party.Get("/swagger/*any", swagger.CustomWrapHandler(&config, swaggerFiles.Handler))
		// swagger配置方法二：默认文档
		party.Get("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	*/

	//config := &swagger.Config{
	//	URL:         "http://localhost:8888/swagger/doc.json", //The url pointing to API definition
	//	DeepLinking: true,
	//}

	//party.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	//party.Get("/swagger/*any", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	party.Get("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
}
