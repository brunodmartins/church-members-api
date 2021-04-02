package main

import (
	"github.com/BrunoDM2943/church-members-api/cmd/gin"
	_ "github.com/BrunoDM2943/church-members-api/internal/infra/config"
	_ "github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
)

func main() {
	gin.StartGinGonicHandler()
}
