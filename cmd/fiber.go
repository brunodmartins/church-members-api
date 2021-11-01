package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/api/middleware"
	"github.com/BrunoDM2943/church-members-api/platform/cdi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"log"
)

//FiberApplication to use as HTTP API
type FiberApplication struct{}

func (FiberApplication) Run() {
	logrus.Info("Init Fiber application")

	app := provideFiberApplication()
	logrus.Info("Application initialized")
	log.Fatal(app.Listen(":8080"))
}

func provideFiberApplication() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ApiErrorMiddleWare,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()
	authHandler := cdi.ProvideAuthHandler()
	userHandler := cdi.ProvideUserHandler()

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("/pong")
	})

	authHandler.SetUpRoutes(app)
	app.Use(middleware.AuthMiddlewareMiddleWare)

	memberHandler.SetUpRoutes(app)
	reportHandler.SetUpRoutes(app)
	userHandler.SetUpRoutes(app)
	return app
}
