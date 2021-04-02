package main

import (
	"github.com/BrunoDM2943/church-members-api/cmd/gin"
	"github.com/BrunoDM2943/church-members-api/internal/infra/config"
	"github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
)

func main() {
	config.BootStrapConfiguration()
	i18n.BootStrapI18N()
	gin.StartGinGonicHandler()
}
