package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Je33/imperial_fleet/internal/config"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql/spaceship"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql/user"
	"github.com/Je33/imperial_fleet/internal/service"
	"github.com/Je33/imperial_fleet/internal/transport/rest/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

// TODO: create swagger doc

func RunRest() error {

	ctx := context.Background()

	cfg := config.Get()

	// connect db
	db, err := mysql.Connect(ctx)
	if err != nil {
		return err
	}

	// Automigrate database
	db.AutoMigrate(&user.User{}, &spaceship.Spaceship{}, &spaceship.SpaceshipArmament{}, &spaceship.SpaceshipArmamentQty{})

	// init repositories
	userRepo := user.NewUserRepo(db)
	spaceshipRepo := spaceship.NewSpaceshipRepo(db)

	// init services
	userService := service.NewUserService(userRepo)
	spaceshipService := service.NewSpaceshipService(spaceshipRepo)

	// init handlers
	userHandler := handler.NewUserHandler(userService)
	spaceshipHandler := handler.NewSpaceshipHandler(spaceshipService)

	// init echo
	e := echo.New()
	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")

	// Auth jwt request
	v1.POST("/auth", userHandler.Auth)
	v1.POST("/register", userHandler.Register)

	// Spaceship
	sg := v1.Group("/spaceships")
	sg.Use(echojwt.JWT([]byte(cfg.JWTSecret)))
	sg.GET("", spaceshipHandler.GetAll)
	sg.GET("/:id", spaceshipHandler.GetById)
	sg.POST("", spaceshipHandler.CreateSpaceship)
	sg.POST("/:id", spaceshipHandler.UpdateSpaceship)
	sg.DELETE("/:id", spaceshipHandler.DeleteSpaceship)

	// Start server
	s := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
