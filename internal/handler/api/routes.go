package api

import (
	"github.com/gofiber/fiber/v2"
)

// Routable defines a way to build a REST Controller routes
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
	// swagger:operation GET /members searchMember
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
	//   type: string
	//   description: The member names
	//   required: false
	// - name: active
	//   in: query
	//   type: string
	//   description: The member status active [true,false]
	//   required: false
	// - name: gender
	//   in: query
	//   type: string
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

	// swagger:operation GET /members/anniversaries getAnniversaries
	//
	// Get anniversaries
	//
	// Returns the birthday and marriage anniversaries for the current week
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: The anniversaries information
	//     schema:
	//       "$ref": "#/definitions/AnniversariesResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/members/anniversaries", handler.lastAnniversaries)
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
	//   type: string
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
	//   type: string
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
	// swagger:operation PUT /members/{id}/contact updateContact
	//
	// Update member - Contact
	//
	// Update the contact for the given member
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The member id
	//   required: true
	// - name: member
	//   in: body
	//   description: The contact to be updated
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ContactRequest"
	// responses:
	//   '200':
	//     description: Contact updated
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/members/:id/contact", handler.updateContact)
	// swagger:operation PUT /members/{id}/address updateAddress
	//
	// Update member - Address
	//
	// Update the address for the given member
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The member id
	//   required: true
	// - name: member
	//   in: body
	//   description: The address to be updated
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/AddressRequest"
	// responses:
	//   '200':
	//     description: Address updated
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/members/:id/address", handler.updateAddress)
	// swagger:operation PUT /members/{id}/person updatePerson
	//
	// Update member - Person
	//
	// Update the person for the given member
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The member id
	//   required: true
	// - name: member
	//   in: body
	//   description: The person to be updated
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdatePersonRequest"
	// responses:
	//   '200':
	//     description: Person updated
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/members/:id/person", handler.updatePerson)

	// swagger:operation PUT /members/{id}/baptism updateBaptism
	//
	// Update member - Baptism
	//
	// Update the baptism info for the given member
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The member id
	//   required: true
	// - name: member
	//   in: body
	//   description: The baptism info to be updated
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdateBaptismRequest"
	// responses:
	//   '200':
	//     description: Baptism updated
	//   '404':
	//     description: Member not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/members/:id/baptism", handler.updateBaptism)
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
	//   type: string
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
	//   type: string
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

func (handler *ParticipantHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation POST /participants postParticipant
	//
	// Create participant
	//
	// Register the receiving participant
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: participant
	//   in: body
	//   description: The participant to be registered
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateParticipantRequest"
	// responses:
	//   '201':
	//     description: Participant registered
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
	app.Post("/participants", handler.postParticipant)

	// swagger:operation GET /participants searchParticipant
	//
	// Search participant
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: name
	//   in: query
	//   type: string
	//   description: The participant name
	//   required: false
	// responses:
	//   '200':
	//     description: Participants found
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/participants", handler.searchParticipant)

	// swagger:operation GET /participants/{id} getParticipant
	//
	// Get participant
	//
	// Returns the participant information
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The participant id
	//   required: true
	// responses:
	//   '200':
	//     description: The participant information
	//     schema:
	//       "$ref": "#/definitions/GetParticipantResponse"
	//   '400':
	//     description: Invalid ID
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '404':
	//     description: Participant not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/participants/:id", handler.getParticipant)

	// swagger:operation PUT /participants/{id} updateParticipant
	//
	// Update participant
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The participant id
	//   required: true
	// - name: participant
	//   in: body
	//   description: The participant to be updated
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateParticipantRequest"
	// responses:
	//   '200':
	//     description: Participant updated
	//   '404':
	//     description: Participant not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Put("/participants/:id", handler.updateParticipant)

	// swagger:operation DELETE /participants/{id} retireParticipant
	//
	// Retire participant
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The participant id
	//   required: true
	// responses:
	//   '200':
	//     description: Participant deleted
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Delete("/participants/:id", handler.retireParticipant)
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

	// swagger:operation GET /users/confirm ConfirmUserEmail
	//
	// Confirms a user email
	//
	// Confirms the email for a given user
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: accessToken
	//   in: query
	//   type: string
	//   description: The access token to be confirmed
	//   required: true
	// responses:
	//   '200':
	//     description: User confirm
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/users/confirm", handler.confirmUserEmail)
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

func (h *ChurchHandler) SetUpRoutes(app *fiber.App) {
	// swagger:operation GET /churches/{id}/statistics getStatistics
	//
	// Get statistics
	//
	// Returns the church statistics information
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The church id
	//   required: true
	// responses:
	//   '200':
	//     description: The church statistics information
	//     schema:
	//       "$ref": "#/definitions/ChurchStatisticsResponse"
	//   '400':
	//     description: Invalid request
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/churches/:id/statistics", h.getStatistics)

	// swagger:operation GET /churches/{id} getChurch
	//
	// Get church
	//
	// Returns the church information
	//
	// ---
	// security:
	// - token: []
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   type: string
	//   description: The church id
	//   required: true
	// responses:
	//   '200':
	//     description: The church information
	//     schema:
	//       "$ref": "#/definitions/GetChurchResponse"
	//   '400':
	//     description: Invalid ID
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '404':
	//     description: Church not found
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	app.Get("/churches/:id", h.getChurchByID)
}
