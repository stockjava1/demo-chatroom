package controller

import (
	"github.com/JabinGP/demo-chatroom/model"
	"github.com/JabinGP/demo-chatroom/model/reqo"
	"github.com/kataras/iris/v12"
)

// Chat with gpt
// @Tags Chatgpt
// @Summary Chat 聊天
// @Description chat with chatgpt
// @Accept  json
// @Produce  json
// Param Authorization header string true "Bearer token"
// @Param message body reqo.PostQuestion true "Account Info"
// @Success 200 {object} reso.PostQuestion
// @Failure 400 {object} reso.HTTPError
// @Failure 404 {object} reso.HTTPError
// @Failure 500 {object} reso.HTTPError
// @Router /chat [post]

func Chat(ctx iris.Context) {
	req := reqo.PostQuestion{}
	ctx.ReadJSON(&req)

	ctx.ContentType("text/html")
	ctx.Header("Transfer-Encoding", "chunked")

	err := chatgptService.Ask(req.Content, true, ctx.ResponseWriter())
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(model.ErrorAskQuestion(err))
		return
	}

	/*
		i := 0
		ints := []int{1, 2, 3, 5, 7, 9, 11, 13, 15, 17, 23, 29}

		ctx.StreamWriter(func(w io.Writer) error {
			fmt.Fprintf(w, "Message number %d<br>", ints[i])

			time.Sleep(500 * time.Millisecond) // simulate delay.
			if i == len(ints)-1 {
				//ctx.Done() //关闭并刷新
				return errors.New("done") //继续写入数据
			}
			i++
			msgLog.Info("test %d", i)
			return nil //继续写入数据
		})
	*/
}
