package controller

import (
	"encoding/json"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/JabinGP/demo-chatroom/model"
	"github.com/JabinGP/demo-chatroom/model/reqo"
	"github.com/JabinGP/demo-chatroom/model/reso"
	"github.com/kataras/iris/v12"
)

var msgLog *logger.CustZeroLogger

func init() {
	msgLog = logger.NewLoggerModule("message")
}

// PostMessage send message
func PostMessage(ctx iris.Context) {
	req := reqo.PostMessage{}
	ctx.ReadJSON(&req)
	xx, err := json.Marshal(ctx.Values().Get("logined").(model.Logined))
	msgLog.Error().Msgf(">>>> post %s", string(xx))
	logined := ctx.Values().Get("logined").(model.Logined)

	insertID, err := messageService.Insert(logined.ID, req.ReceiverID, req.Content)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(model.ErrorInsertDatabase(err))
		return
	}

	res := reso.PostMessage{
		ID: insertID,
	}

	ctx.JSON(res)
}

// GetMessage get all message
func GetMessage(ctx iris.Context) {
	req := reqo.GetMessage{}
	ctx.ReadQuery(&req)
	logined := ctx.Values().Get("logined").(model.Logined)

	msgList, err := messageService.Query(
		req.BeginID,
		req.BeginTime,
		req.EndTime,
		logined.ID,
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(model.ErrorQueryDatabase(err))
		return
	}

	// Build response object
	resList := []reso.GetMessage{}

	for _, msg := range msgList {
		private := false
		if len(msg.Receiver.ID) > 0 {
			private = true
		}
		// Get single res
		res := reso.GetMessage{
			ID:         msg.Message.ID,
			SenderID:   msg.Message.SenderID,
			SenderName: msg.Sender.Username,
			Content:    msg.Message.Content,
			SendTime:   msg.Message.SendTime,
			Private:    private,
		}

		// Add into resList
		resList = append(resList, res)
	}

	ctx.JSON(resList)
}
