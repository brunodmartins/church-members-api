// church-members-api
//
// This API manages the members of a church.
//
//     Schemes: http
//     Host: localhost:8080
//     Version: 1.0.0
//     basePath: /
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//	   - application/csv
//	   - application/pdf
//
// swagger:meta
package main

import (
	"github.com/BrunoDM2943/church-members-api/cmd"
	config "github.com/BrunoDM2943/church-members-api/platform/config"
	i18n "github.com/BrunoDM2943/church-members-api/platform/i18n"
)

//go:generate swagger generate spec -m -o ./docs/specs/swagger.yaml
func main() {
	config.InitConfiguration()
	_ = i18n.GetMessageService()
	cmd.ProvideRunner().Run()
}
