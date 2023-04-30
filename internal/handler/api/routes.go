package api

import (
	"github.com/gofiber/fiber/v2"
)

// Routable defines a way to builds a REST Controller routes
type Routable interface {
	//SetUpRoutes build the REST controller routes
	SetUpRoutes(app *fiber.App)
}

func (handler *MemberHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation POST /members postMember
	//
	// Create member
	//
	// Register the receiving member
	//
	// ---
	// security:
	// - token: []
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
	// swagger:operation POST /members/search searchMember
	//
	// Search member
	//
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: name
	//   in: query
	//   description: The member names
	//   required: false
	// - name: active
	//   in: query
	//   description: The member status active [true,false]
	//   required: false
	// - name: gender
	//   in: query
	//   description: The member gender [M,F]
	//   required: false
	// responses:
	//   '200':
	//     description: Members found
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/members", handler.searchMember)
	// swagger:operation GET /members/{id} getMember
	//
	// Get member
	//
	// Returns the member information
	//
	// ---
	// security:
	// - token: []
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
	// swagger:operation DELETE /members/{id} retireMember
	//
	// Retire Member
	//
	// Retire a member from the church
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: The member id
	//   required: true
	// - name: body
	//   in: body
	//   description: The retire information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RetireMemberRequest"
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
	app.Delete("/members/:id", handler.retireMember)
}

func (handler *ReportHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation POST /reports/members/birthday generateBirthDayReport
	//
	// Birthday report
	//
	// Generates a CSV birthday report
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/csv
	// responses:
	//   '200':
	//     description: CSV report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/reports/members/birthday", handler.generateBirthDayReport)
	// swagger:operation POST /reports/members/marriage generateMarriageReport
	//
	// Marriage report
	//
	// Generates a CSV Marriage report
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/csv
	// responses:
	//   '200':
	//     description: CSV report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/reports/members/marriage", handler.generateMarriageReport)
	// swagger:operation POST /reports/members/legal generateLegalReport
	//
	// Legal report
	//
	// Generates a PDF legal report
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/pdf
	// responses:
	//   '200':
	//     description: PDF report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/reports/members/legal", handler.generateLegalReport)
	// swagger:operation POST /reports/members/classification/{classification} generateClassificationReport
	//
	// Member report
	//
	// Generates a PDF member report by classification
	//
	// ---
	// security:
	// - token: []
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
	app.Post("/reports/members/classification/:classification", handler.generateClassificationReport)
	// swagger:operation POST /reports/members generateMembersReport
	//
	// Member report
	//
	// Generates a PDF member report
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/pdf
	// responses:
	//   '200':
	//     description: PDF report
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/reports/members", handler.generateMembersReport)

	// swagger:operation GET /reports/{reportType} getURLForReport
	//
	// Get a report file
	//
	// Returns a report file url
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: reportType
	//   in: path
	//   description: The report type [members,legal,classification,birthdate,marriage]
	//   required: true
	// responses:
	//   '307':
	//     description: S3 url to redirect
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/reports/:reportType", handler.getURLForReport)
}

func (handler *AuthHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation GET /users/token getToken
	//
	// Get a user Token
	//
	// Generates a token for a given user
	//
	// ---
	// security:
	// - login: []
	// produces:
	// - application/csv
	// responses:
	//   '201':
	//     description: Token generated
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/users/token", handler.getToken)
}

func (handler *UserHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation POST /users postUser
	//
	// Create a user
	//
	// Register the receiving user
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: user
	//   in: body
	//   description: The user to be registered
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateUserRequest"
	// responses:
	//   '201':
	//     description: User registered
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Post("/users", handler.PostUser)
}
