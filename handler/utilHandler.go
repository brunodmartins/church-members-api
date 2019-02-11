package handler

import (
	"time"

	"github.com/BrunoDM2943/church-members-api/utils"
	"github.com/gin-gonic/gin"
)

type UtilHandler struct {
	service utils.UtilsService
}

func (handler *UtilHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/utils/aniversariantesMes", handler.GetBirthDayMembers)
}

func NewUtilHandler(service utils.UtilsService) *UtilHandler {
	return &UtilHandler{
		service: service,
	}
}

func (handler *UtilHandler) GetBirthDayMembers(c *gin.Context) {
	date := time.Now()

	list, err := handler.service.FindMonthBirthday(date)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, list)
	return
}
