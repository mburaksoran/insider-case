package main

import (
	"context"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	apiHandler "github.com/mburaksoran/insider-case/internal/app/api"
	backgroundjob "github.com/mburaksoran/insider-case/internal/app/background_job"
	"github.com/mburaksoran/insider-case/internal/app/background_job/handler"
	"github.com/mburaksoran/insider-case/internal/app/background_job/tasks"
	"github.com/mburaksoran/insider-case/internal/app/config"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/service"
	"github.com/mburaksoran/insider-case/internal/infra/clients"
	"github.com/mburaksoran/insider-case/internal/infra/engines"
	"github.com/mburaksoran/insider-case/internal/infra/repository"
	zaploki "github.com/paul-milne/zap-loki"
	"github.com/pressly/goose"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"

	"github.com/mburaksoran/insider-case/internal/infra/workerpool"

	_ "github.com/mburaksoran/insider-case/docs"
)

// @title			Insider Case Study
// @version		1.0
// @description	This is a sample server for Insider Case Study.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	support@insider.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
func main() {
	logger := prepareLogger()

	vaultEngine, err := engines.SetVaultEngine()
	if err != nil {
		logger.Fatalf("Vault client creation error: %s", err)
	}

	vaultSecrets, err := vaultEngine.GetSecret()
	if err != nil {
		logger.Fatalf("Error reading secret: %s", err)
	}

	var appConfig = &config.AppConfig{}
	err = appConfig.MapVaultSecretToConfig(vaultSecrets)
	if err != nil {
		logger.Fatalf("Error reading secret: %s", err)
	}

	sqlEngine, err := engines.SetSqlDBEngine(appConfig)
	if err != nil {
		log.Fatalf("Error creating sqlEngine: %s", err)
	}

	if err := goose.Up(sqlEngine.Client, "./internal/infra/migrations"); err != nil {
		log.Fatalf("Migration işlemi başarısız: %v", err)
	}

	err = engines.SetRedisEngine(appConfig, logger)
	if err != nil {
		log.Fatalf("Error creating redisEngine: %s", err)
	}
	ch := make(chan *models.BackgroundJob)
	errorChannel := make(chan error, 1)

	bgRepo := repository.NewBackgroundJobRepository(sqlEngine.Client)
	messageRepo := repository.NewMessageRepository(sqlEngine.Client)
	redisMessageRepo := repository.NewRedisMessageRepository(logger)

	httpClient := clients.NewHttpClient(appConfig, logger)
	bgService := service.NewBackgroundJobService(bgRepo)
	messageService := service.NewMessageService(messageRepo)
	redisMessageService := service.NewRedisMessageService(redisMessageRepo, logger)

	messageHandler := apiHandler.NewMessageHandler(messageService, logger)
	backgroundJobHandler := apiHandler.NewBackgroundJobHandler(bgService, logger)

	factory := handler.NewHandlerFactory(bgService, messageService, redisMessageService, httpClient, logger)

	fetcher := backgroundjob.NewFetcher(sqlEngine.Client, bgService, ch, errorChannel, logger)

	go fetcher.Start()

	app := fiber.New()

	prometheus := fiberprometheus.New("insider-case-study")
	prometheus.RegisterAt(app, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping"})
	app.Use(prometheus.Middleware)

	prepareRoutes(app, messageHandler, backgroundJobHandler)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			logger.Fatal(err)
		}
	}()

	pool, _ := workerpool.NewWorkerPool(5, 5, logger)

	pool.Start()
	dispatcher := backgroundjob.NewDispatcher(errorChannel, &config.AppConfig{}, pool)

	defer pool.Stop()

	go receiveMessage(ch, dispatcher, factory, errorChannel)
	select {}
}
func prepareLogger() *zap.SugaredLogger {
	zapConfig := zap.NewProductionConfig()

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	loki := zaploki.New(context.Background(), zaploki.Config{
		Url:          "http://loki:3100",
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": "insider-case-study"},
	})

	logger, err := loki.WithCreateLogger(zapConfig)
	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}

func receiveMessage(ch chan *models.BackgroundJob, dispatcher *backgroundjob.Dispatcher, factory *handler.HandlerFactory, errorChannel chan error) {
	for msg := range ch {
		dispatcher.DispatchJob(&tasks.JobTask{
			Job:     msg,
			Factory: factory,
			Ctx:     context.Background(),
			ErrorCh: errorChannel,
		})
	}
}

func prepareRoutes(app *fiber.App, msgHandler apiHandler.MessageHandlerInterface, bgHandler apiHandler.BackgroundJobHandlerInterface) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	// HealthcheckHandler godoc
	//	@Summary		Healthcheck API
	//	@Description	Get the health status of the application
	//	@Tags			Health
	//	@Accept			json
	//	@Produce		json
	//	@Router			/health-check [get]
	app.Get("/health-check", apiHandler.HealthcheckHandler)

	// HealthAliveHandler godoc
	//	@Summary		Health Alive API
	//	@Description	Check if the application is alive
	//	@Tags			Health
	//	@Accept			json
	//	@Produce		json
	//	@Router			/health-alive [get]
	app.Get("/health-alive", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"health": "ok"})
	})

	app.Put("/jobs/start", bgHandler.StartJobsHandler)

	app.Put("/jobs/stop", bgHandler.StopJobsHandler)

	app.Get("/messages/sent", msgHandler.ListSendMessagesHandler)

}
