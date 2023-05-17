package route

import (
	"github.com/JabinGP/demo-chatroom/controller"
	"github.com/JabinGP/demo-chatroom/middleware"
	"github.com/kataras/iris/v12/core/router"
)

func routeUser(party router.Party) {
	party.Post("/login", controller.PostLogin)

	// @Summary GetUser return user list
	// @Description
	// @Tags　学生
	// @Accept application/json
	// @Produce application/json
	// Param Authorization header string true "Bearer token"
	// Param object query models.pojo.User false "查询参数"
	// @Param name query string true "name"
	// Param body st.PaySearch true "交款查询参数"
	// @Param tel query string true "tel"
	// Success 200 {object} models.Result
	// Failure 400 {object} models.Result
	// @Router /user [get/post]
	party.Post("/user", controller.PostUser)
	party.Get("/user", controller.GetUser)
	party.Put("/user", middleware.JWT.Serve, middleware.Logined, controller.PutUser)
}
