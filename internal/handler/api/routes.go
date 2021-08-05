package api

import (
	"github.com/gofiber/fiber/v2"
)

//Routable defines a way to builds a REST Controller routes
type Routable interface {
	//SetUpRoutes build the REST controller routes
	SetUpRoutes(app *fiber.App)
}

func (handler *MemberHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation POST members postMember
	//
	// Create member
	//
	// Register the receiving member
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: member
	//   in: body
	//   description: The member to be registered
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateMemberRequest"
	// responses:
	//   '201':
	//     description: Member registered
	//     schema:
	//       "$ref": "#/definitions/CreateMemberResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/members", handler.postMember)
	// swagger:operation POST members/search searchMember
	//
	// Search member
	//
	// A GraphQL endpoint to search for members data
	// {
	//		member(name:"Bruno", active:true, gender:"M"){
	//			  id
	//				person{
	//					firstName,
	//					lastName
	//				}
	//		}
	// }
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: query
	//   in: body
	//   description: The GraphQL query
	//   required: true
	// responses:
	//   '200':
	//     description: Members found
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/GraphQLErrorResponse"
	app.Post("/members/search", handler.searchMember)
	// swagger:operation GET members/{id} getMember
	//
	// Get member
	//
	// Returns the member information
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: The member id
	//   required: true
	// responses:
	//   '200':
	//     description: The member information
	//     schema:
	//       "$ref": "#/definitions/GetMemberResponse"
	//   '400':
	//     description: Invalid ID
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/members/:id", handler.getMember)
	// swagger:operation PUT members/{id}/status putMemberStatus
	//
	// Put member status
	//
	// Changes the member status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: The member id
	//   required: true
	// - name: body
	//   in: body
	//   description: The status information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/PutMemberStatusRequest"
	// responses:
	//   '200':
	//     description: Status change successfully
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/members/:id/status", handler.putStatus)
}

func (handler *ReportHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation GET reports/members/birthday generateBirthDayReport
	//
	// Birthday report
	//
	// Generates a CSV birthday report
	//
	// ---
	// produces:
	// - application/csv
	// responses:
	//   '200':
	//     description: CSV report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/members/birthday", handler.generateBirthDayReport)
	// swagger:operation GET reports/members/marriage generateMarriageReport
	//
	// Marriage report
	//
	// Generates a CSV Marriage report
	//
	// ---
	// produces:
	// - application/csv
	// responses:
	//   '200':
	//     description: CSV report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/members/marriage", handler.generateMarriageReport)
	// swagger:operation GET reports/members/legal generateLegalReport
	//
	// Legal report
	//
	// Generates a PDF legal report
	//
	// ---
	// produces:
	// - application/pdf
	// responses:
	//   '200':
	//     description: PDF report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/members/legal", handler.generateLegalReport)
	// swagger:operation GET reports/members/classification/{classification} generateClassificationReport
	//
	// Member report
	//
	// Generates a PDF member report by classification
	//
	// ---
	// produces:
	// - application/pdf
	// parameters:
	// - name: classification
	//   in: path
	//   description: The member classification [adult, teen, young, children]
	//   required: true
	// responses:
	//   '200':
	//     description: PDF report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/members/classification/:classification", handler.generateClassificationReport)
	// swagger:operation GET reports/members generateMembersReport
	//
	// Member report
	//
	// Generates a PDF member report
	//
	// ---
	// produces:
	// - application/pdf
	// responses:
	//   '200':
	//     description: PDF report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/members", handler.generateMembersReport)
}
