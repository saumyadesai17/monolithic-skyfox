package server

import (
	"fmt"
	"net/http"
	"time"

	"skyfox/bookings/constants"
	"skyfox/bookings/controller"
	"skyfox/bookings/database/connection"
	database "skyfox/bookings/database/seed"
	"skyfox/bookings/repository"
	"skyfox/bookings/service"
	"skyfox/common/logger"
	"skyfox/common/middleware/cors"
	"skyfox/common/middleware/security"
	"skyfox/common/middleware/validator"
	appConf "skyfox/config"
	movieservice "skyfox/movieservice/movie_gateway"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	_ "skyfox/docs" //indirect

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(cfg appConf.AppConfig) error {

	logger.InitAppLogger(cfg.Logger)

	handler := connection.NewDBHandler(cfg.Database)
	db := handler.Instance()

	movieGateway := movieservice.NewMovieGateway(cfg.MovieGateway)

	// instantiate repositories
	bookingRepository := repository.NewBookingRepository(db)
	showRepository := repository.NewShowRepository(db)
	userRepository := repository.NewUserRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	userAccountRepository := repository.NewAccountRepository(db)

	database.SeedDB(userRepository)

	// instantiate all services
	bookingService := service.NewBookingService(bookingRepository, showRepository)
	bookingService.SetCustomerRepository(customerRepository)
	showService := service.NewShowService(showRepository, movieGateway)
	userService := service.NewUserService(userRepository)
	revenueService := service.NewRevenueService(bookingRepository, showRepository)
	authService := service.NewAuthService(userAccountRepository)

	// instantiate all handlers
	bookingController := controller.NewBookingController(bookingService)
	showController := controller.NewShowController(showService)
	userController := controller.NewUserController(userService)
	revenueController := controller.NewRevenueController(revenueService)
	authController := controller.NewAuthController(authService)

	router := setupApp(cfg)

	authRouter := routerGroupWithAuth(router, userService)
	noAuthRouter := routerGroupWithNoAuth(router)

	booking := authRouter.Group(constants.BookingEndPoint)
	{
		booking.POST("", bookingController.CreateBooking)
	}

	revenue := authRouter.Group(constants.RevenueEndPoint)
	{
		revenue.GET("", revenueController.GetRevenue)
	}

	show := authRouter.Group(constants.ShowEndPoint)
	{
		show.GET("", showController.Shows)
	}

	authRouter.GET(constants.LoginEndPoint, userController.Login)

	// Versioned public API routes
	v1 := noAuthRouter.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authController.Signup)
		}
	}

	noAuthRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := start(router, cfg.Server)
	if err != nil {
		return err
	}
	return nil
}

func start(r *gin.Engine, cfg appConf.ServerConfig) error {
	s := &http.Server{
		Addr:         port(cfg),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		return fmt.Errorf("unable to start gin server. error: %w", err)
	}
	return nil
}

func setupApp(cfg appConf.AppConfig) *gin.Engine {
	gin.SetMode(cfg.Server.GineMode)
	engine := gin.New()
	binding.Validator = new(validator.DtoValidator)
	return setupMiddleware(engine, cfg)
}

func setupMiddleware(engine *gin.Engine, cfg appConf.AppConfig) *gin.Engine {
	engine.Use(cors.SetupCORS())
	engine.Use(ginzap.Ginzap(logger.GetLogger(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger.GetLogger(), true))
	return engine
}

func setAuthMiddleware(rg *gin.RouterGroup, us controller.UserService) *gin.RouterGroup {
	rg.Use(security.Authenticate(us))
	return rg
}

func port(c appConf.ServerConfig) string {
	return fmt.Sprintf(":%d", c.Port)
}

func routerGroupWithAuth(engine *gin.Engine, userService controller.UserService) *gin.RouterGroup {
	authRouter := engine.Group("")
	authRouter = setAuthMiddleware(authRouter, userService)
	return authRouter
}

func routerGroupWithNoAuth(engine *gin.Engine) *gin.RouterGroup {
	noAuthRouter := engine.Group("")
	return noAuthRouter
}
