// church-members-api
//
// This API manages the members of a church.
//
//	Schemes: http
//	Host: localhost:8080
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	   - application/csv
//	   - application/pdf
//
//	Security:
//	- login:
//	- token:
//
//	SecurityDefinitions:
//	  login:
//	    type: basic
//	  token:
//	    type: apiKey
//	    name: X-Auth-Token
//	    in: header
//
// swagger:meta
package main

import (
	"github.com/brunodmartins/church-members-api/cmd"
	"github.com/brunodmartins/church-members-api/platform/config"
)

//go:generate swagger generate spec -m -o ./docs/specs/swagger.yaml
func main() {
	config.InitConfiguration()
	cmd.ProvideRunner().Run()
}
