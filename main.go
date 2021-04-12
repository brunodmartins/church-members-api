package main

import (
	"github.com/BrunoDM2943/church-members-api/cmd"
	_ "github.com/BrunoDM2943/church-members-api/internal/infra/config"
	_ "github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
)

func main() {
	cmd.ProvideRunner().Run()
}
