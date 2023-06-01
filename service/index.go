package service

import (
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/database"
	"github.com/JabinGP/demo-chatroom/infra/logger"
)

// NewMessage get a message service
func NewMessage() MessageService {
	return MessageService{
		db: database.DB,
	}
}

// NewUser get a user service
func NewUser() UserService {
	return UserService{
		db: database.DB,
	}
}

// NewUser get a user service
func NewChatgpt() ChatGptService {
	return ChatGptService{config.Viper.GetString("openai.systemUser"), logger.NewLoggerModule("chatgpt")}
}
