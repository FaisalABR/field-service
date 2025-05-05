package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type ITimeRoute interface {
	Run()
}

func NewTimeRoute(controler controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) ITimeRoute {
	return &TimeRoute{
		controller: controler,
		group:      group,
		client:     client,
	}
}

func (t *TimeRoute) Run() {
	group := t.group.Group("/time")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.
		CheckRole([]string{constants.Admin}, t.client),
		t.controller.GetTime().GetAll)
	group.POST("", middlewares.
		CheckRole([]string{constants.Admin}, t.client),
		t.controller.GetTime().Create)
	group.GET("/:uuid", middlewares.
		CheckRole([]string{constants.Admin}, t.client),
		t.controller.GetTime().GetByUUID)
}
