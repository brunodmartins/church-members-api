package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/BrunoDM2943/church-members-api/internal/infra/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type LambdaApplication struct {}

func (LambdaApplication) Run(){
	router := provideGinGonic()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()

	memberHandler.SetUpRoutes(router)
	reportHandler.SetUpRoutes(router)

	router.Use(func(c *gin.Context) {
		c.Next()
		headers := c.Writer.Header()
		type responseBodyWriter struct {
			gin.ResponseWriter
			body *bytes.Buffer
		}
		if headers.Get("content-type") == "application/pdf" && c.Writer.Status() == 200 {
			file, _ := c.GetRawData()
			newBody , _ := json.Marshal(gin.H{
				"body": base64.RawStdEncoding.EncodeToString(file),
				"isBase64Encoded": true,
			})
			c.Writer = &responseBodyWriter{body: bytes.NewBuffer(newBody), ResponseWriter: c.Writer}
		}
	})

	ginLambda := ginadapter.New(router)

	lambda.Start(ginLambda.ProxyWithContext)
}