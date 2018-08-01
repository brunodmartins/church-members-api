package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	SetUpRoutes(routes *gin.Engine)
}
