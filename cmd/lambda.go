package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/infra/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

//LambdaApplication implements the Application interface to execute this app on AWS Lambda
type LambdaApplication struct{}

//Run starts a lambda adapter on top of gin-gonic to execute the application on serverless
func (LambdaApplication) Run() {
	router := provideGinGonic()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()

	memberHandler.SetUpRoutes(router)
	reportHandler.SetUpRoutes(router)

	ginLambda := ginadapter.New(router)

	lambda.Start(ginLambda.ProxyWithContext)
}
