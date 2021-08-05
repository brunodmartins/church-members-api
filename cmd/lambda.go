package cmd

import (
	"github.com/BrunoDM2943/church-members-api/platform/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
)

//LambdaApplication implements the Application interface to execute this app on AWS Lambda
type LambdaApplication struct{}

//Run starts a lambda adapter on top of gin-gonic to execute the application on serverless
func (LambdaApplication) Run() {
	app := provideFiberApplication()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()

	memberHandler.SetUpRoutes(app)
	reportHandler.SetUpRoutes(app)

	fiberLambda := fiberadapter.New(app)

	lambda.Start(fiberLambda.ProxyWithContext)
}
