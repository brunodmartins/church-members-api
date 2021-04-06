package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/infra/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

type LambdaApplication struct {}

func (LambdaApplication) Run(){
	router := provideGinGonic()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()

	memberHandler.SetUpRoutes(router)
	reportHandler.SetUpRoutes(router)

	ginLambda := ginadapter.New(router)

	lambda.Start(ginLambda.ProxyWithContext)
}