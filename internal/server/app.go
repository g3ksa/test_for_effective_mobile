package server

import (
	"UserService/internal/users"
	"UserService/internal/users/delivery"
	"UserService/internal/users/repository/psql"
	"UserService/internal/users/service"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type App struct {
	httpServer  *http.Server
	userService users.Service
}

func NewApp(log *slog.Logger) *App {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName))
	if err != nil {
		log.Error("database error", slog.String("error", err.Error()))
	}
	userRepo := psql.New(db, log)

	userService := service.New(
		userRepo,
		log,
	)
	return &App{
		userService: userService,
	}
}

func (a *App) Run(port string, log *slog.Logger) error {
	router := gin.Default()

	router.Use(gin.Recovery(), gin.Logger())

	api := router.Group("/")
	delivery.RegisterHTTPEndpoints(api, a.userService, log)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	slog.Info("Server started on port " + port)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Error("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
