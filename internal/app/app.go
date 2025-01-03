package app

import (
	"context"
	"errors"
	"github.com/1Storm3/flibox-api/internal/shared/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/zap"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/controller/http"
	"github.com/1Storm3/flibox-api/internal/delivery/grpc"
	"github.com/1Storm3/flibox-api/internal/delivery/middleware"
	"github.com/1Storm3/flibox-api/internal/delivery/rest"
	"github.com/1Storm3/flibox-api/internal/metrics"
	"github.com/1Storm3/flibox-api/internal/metrics/interceptor"
	"github.com/1Storm3/flibox-api/internal/repo"
	"github.com/1Storm3/flibox-api/internal/service"
	"github.com/1Storm3/flibox-api/pkg/kafka"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type App struct {
	httpServer *fiber.App
}

func New() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {

	a.initFiberServer()

	a.initCORS()

	err := a.initMetrics(ctx)
	if err != nil {
		return err
	}

	cfg := config.MustLoad()

	storage, err := postgres.NewStorage(cfg.DB.URL)
	if err != nil {
		logger.Error("Error while creating storage: %v", zap.Error(err))
		return err
	}
	defer storage.Close()

	grpcClient, err := grpc.NewClient(cfg)
	if err != nil {
		logger.Error("Error while creating grpc client: %v", zap.Error(err))
		return err
	}
	defer grpcClient.Close()

	// s3
	s3Service, err := service.NewS3Service(cfg)
	if err != nil {
		logger.Error("Error while creating s3 service: %v", zap.Error(err))
	}
	// email
	emailService := service.NewEmailService(cfg)

	//kafka
	kafkaProducer := kafka.NewProducer([]string{"host.docker.internal:9092"}, "favorite")
	// token
	tokenService := service.NewTokenService()
	// user
	userRepo := repo.NewUserRepo(storage)
	userService := service.NewUserService(userRepo, s3Service)
	userController := http.NewUserController(userService)

	// auth
	authService := service.NewAuthService(userService, emailService, cfg, tokenService, helper.TakeHTMLTemplate)
	authController := http.NewAuthController(authService)

	// external
	externalService := service.NewExternalService(cfg)
	externalController := http.NewExternalController(externalService, s3Service)
	// film
	filmRepo := repo.NewFilmRepo(storage)
	filmService := service.NewFilmService(filmRepo, externalService)
	filmController := http.NewFilmController(filmService)

	// film sequel
	filmSequelRepo := repo.NewFilmSequelRepo(storage)
	filmSequelService := service.NewFilmSequelService(filmSequelRepo, filmService, cfg)
	filmSequelController := http.NewFilmSequelController(filmSequelService)

	// film similar
	filmSimilarRepo := repo.NewFilmSimilarRepo(storage)
	filmSimilarService := service.NewFilmSimilarService(filmSimilarRepo, filmService, cfg)
	filmSimilarController := http.NewFilmSimilarController(filmSimilarService)

	// user film
	userFilmRepo := repo.NewUserFilmRepo(storage)
	userFilmService := service.NewUserFilmService(userFilmRepo)

	// history films
	historyFilmsRepo := repo.NewHistoryFilmsRepo(storage)
	historyFilmsService := service.NewHistoryFilmsService(historyFilmsRepo)

	// recommend
	recommendService := service.NewRecommendService(historyFilmsService, filmService, userFilmService, grpcClient)

	userFilmController := http.NewUserFilmController(userFilmService, filmService, recommendService, kafkaProducer)
	historyFilmsController := http.NewHistoryFilmsController(historyFilmsService, recommendService)

	// comment
	commentRepo := repo.NewCommentRepo(storage)
	commentService := service.NewCommentService(commentRepo)
	commentController := http.NewCommentController(commentService)

	// collection
	collectionRepo := repo.NewCollectionRepo(storage)
	collectionService := service.NewCollectionService(collectionRepo)
	collectionController := http.NewCollectionController(collectionService)

	// collection film
	collectionFilmRepo := repo.NewCollectionFilmRepo(storage)
	collectionFilmService := service.NewCollectionFilmService(collectionFilmRepo)
	collectionFilmController := http.NewCollectionFilmController(collectionFilmService)

	authMiddleware := middleware.AuthMiddleware(userRepo, cfg, tokenService)

	router := rest.NewRouter(filmController,
		filmSequelController,
		userController,
		filmSimilarController,
		userFilmController,
		authController,
		externalController,
		commentController,
		collectionController,
		collectionFilmController,
		historyFilmsController,
	)
	router.LoadRoutes(a.httpServer, authMiddleware)

	a.httpServer.Get("/swagger/*", fiberSwagger.WrapHandler)
	a.httpServer.Use(interceptor.MetricsInterceptor())

	go func() {
		if err := a.httpServer.Listen(":8080"); err != nil && !errors.Is(err, fiber.ErrBadGateway) {
			logger.Fatal("Error while starting server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	return a.httpServer.Shutdown()
}

func (a *App) initFiberServer() {
	a.httpServer = fiber.New(fiber.Config{
		ErrorHandler: a.customErrorHandler(),
	})
}

func (a *App) initCORS() {
	a.httpServer.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
}

func (a *App) initMetrics(ctx context.Context) error {
	return metrics.Init(ctx)
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.MustLoad().Env)
}

func (a *App) customErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var message string

		var httpErr *sys.Error
		if errors.As(err, &httpErr) {
			code = sys.ErrorMap[httpErr.Message]
			message = httpErr.Error()
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}

		return ctx.Status(code).JSON(fiber.Map{
			"statusCode": code,
			"message":    message,
		})
	}
}
