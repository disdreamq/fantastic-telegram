// @title           Blog API
// @version         1.0
// @description     REST API для блога с авторизацией, управлением пользователями и постами
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in                        header
// @name                      Authorization
// @description               Введите токен в формате: Bearer <ваш_токен>

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disdreamq/fantastic-telegram/services/post/config"
	"github.com/disdreamq/fantastic-telegram/services/post/internal/repository/postgres"
	"github.com/disdreamq/fantastic-telegram/services/post/internal/repository/redis"
	"github.com/disdreamq/fantastic-telegram/services/post/internal/service"

	_ "github.com/disdreamq/fantastic-telegram/services/post/docs"

	"github.com/disdreamq/fantastic-telegram/services/post/internal/handler"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	// Logging
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger.Info().Msg("Starting the application.")

	// Load cfg
	if _, err := os.Stat("../../.env"); err == nil {
		err := godotenv.Load("../../.env")
		if err != nil {
			logger.Fatal().Err(err).Msg("Fatal error during parse env file.")
		}
	}
	cfg := config.Load()

	// Connect to redis
	rdb, err := redis.RedisConnect(cfg)
	if err != nil {
		logger.Err(err).Str("component", "Redis").Msg("Redis could not connect to db.")
	}
	cache := redis.NewRedisCache(rdb)

	// Connect to DB
	DB, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Str("component", "Postgres").Msg("Postgres could not connect to db.")
	}
	// Prepare post controller
	postRepo := postgres.NewPostRepository(DB)
	postSVC := service.NewPostService(postRepo, cache)
	postCtrl := handler.NewPostController(postSVC)

	r := handler.NewRouter(rdb, postCtrl, cfg.ProtectedRPM, cfg.PublicRPM, logger)

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Start server
	srv := &http.Server{
		Addr:    ":" + strconv.FormatInt(int64(cfg.HttpPort), 10),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().
				Err(err).
				Msg("Critical error during starting server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	// Shutdown
	logger.Info().Msg("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().
			Err(err).
			Msg("Server shutdown failed")
	}
	DB.Close()
	rdb.Close()

}
