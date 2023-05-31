package controller

import (
	"github.com/JabinGP/demo-chatroom/infra/database"
)

import "github.com/JabinGP/demo-chatroom/service"

var db = database.DB
var messageService = service.NewMessage()
var userService = service.NewUser()

var chatgptService = service.NewChatgpt()
