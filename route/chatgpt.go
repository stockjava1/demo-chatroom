package route

import (
	"github.com/JabinGP/demo-chatroom/controller"
	"github.com/kataras/iris/v12/core/router"
)

func routeChatGpt(party router.Party) {
	party.Post("/chat", controller.Chat)
}
