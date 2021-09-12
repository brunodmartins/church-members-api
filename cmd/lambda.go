package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/api/middleware"
	"github.com/BrunoDM2943/church-members-api/platform/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

//LambdaApplication implements the Application interface to execute this app on AWS Lambda
type LambdaApplication struct{}

//Run starts a lambda adapter on top of gin-gonic to execute the application on serverless
func (LambdaApplication) Run() {
	app := provideFiberApplication()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()
	authHandler := cdi.ProvideAuthHandler()

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("/pong")
	})

	authHandler.SetUpRoutes(app)
	app.Use(middleware.AuthMiddlewareMiddleWare)

	memberHandler.SetUpRoutes(app)
	reportHandler.SetUpRoutes(app)

	fiberLambda := fiberadapter.New(app)

	lambda.Start(fiberLambda.ProxyWithContext)
}
