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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

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

	// User
	//r := v1.Group("/user")
	//r.Use(echojwt.JWT([]byte(cfg.JWTSecret)))
	//r.GET("/me", userHandler.Me)

	// Spaceship
	v1.GET("/spaceships", spaceshipHandler.GetAll)
	v1.GET("/spaceships/:id", spaceshipHandler.GetById)
	v1.POST("/spaceships", spaceshipHandler.CreateSpaceship)
	v1.POST("/spaceships/:id", spaceshipHandler.UpdateSpaceship)
	v1.DELETE("/spaceships/:id", spaceshipHandler.DeleteSpaceship)

	// Start server
	s := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
