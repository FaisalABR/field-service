package controllers

import (
	fieldController "field-service/controllers/field"
	fieldScheduleController "field-service/controllers/fieldschedule"
	timeController "field-service/controllers/time"
	"field-service/services"
)

type ControllerRegistry struct {
	services services.IServiceRegistry
}

type IControllerRegistry interface {
	GetField() fieldController.IFieldController
	GetFieldSchedule() fieldScheduleController.IFieldScheduleController
	GetTime() timeController.ITimeController
}

func NewControllerRegistry(services services.IServiceRegistry) IControllerRegistry {
	return &ControllerRegistry{services: services}
}

func (c *ControllerRegistry) GetField() fieldController.IFieldController {
	return fieldController.NewFieldController(c.services)
}

func (c *ControllerRegistry) GetFieldSchedule() fieldScheduleController.IFieldScheduleController {
	return fieldScheduleController.NewFieldScheduleController(c.services)
}

func (c *ControllerRegistry) GetTime() timeController.ITimeController {
	return timeController.NewTimeController(c.services)
}
