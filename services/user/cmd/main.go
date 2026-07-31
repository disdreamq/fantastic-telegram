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

	"github.com/disdreamq/fantastic-telegram/services/user/config"
	_ "github.com/disdreamq/fantastic-telegram/services/user/docs"
	"github.com/disdreamq/fantastic-telegram/services/user/internal/infra/hasher"
	"github.com/disdreamq/fantastic-telegram/services/user/internal/infra/jwt"
	"github.com/disdreamq/fantastic-telegram/services/user/internal/repository/postgres"
	"github.com/disdreamq/fantastic-telegram/services/user/internal/repository/redis"

	"github.com/disdreamq/fantastic-telegram/services/user/internal/service"
	"github.com/disdreamq/fantastic-telegram/services/user/internal/transport/http/handler"
	"github.com/fatih/color"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

func printBanner(cfg *config.Config) {
	color.NoColor = false // Force colors

	banner := ` ______   ______     __   __     ______   ______     ______     ______   __     ______     ______   ______     __         ______     ______     ______     ______     __    __    
/\  ___\ /\  __ \   /\ "-.\ \   /\__  _\ /\  __ \   /\  ___\   /\__  _\ /\ \   /\  ___\   /\__  _\ /\  ___\   /\ \       /\  ___\   /\  ___\   /\  == \   /\  __ \   /\ "-./  \   
\ \  __\ \ \  __ \  \ \ \-.  \  \/_/\ \/ \ \  __ \  \ \___  \  \/_/\ \/ \ \ \  \ \ \____  \/_/\ \/ \ \  __\   \ \ \____  \ \  __\   \ \ \__ \  \ \  __<   \ \  __ \  \ \ \-./\ \  
 \ \_\    \ \_\ \_\  \ \_\\"\_\    \ \_\  \ \_\ \_\  \/\_____\    \ \_\  \ \_\  \ \_____\    \ \_\  \ \_____\  \ \_____\  \ \_____\  \ \_____\  \ \_\ \_\  \ \_\ \_\  \ \_\ \ \_\ 
  \/_/     \/_/\/_/   \/_/ \/_/     \/_/   \/_/\/_/   \/_____/     \/_/   \/_/   \/_____/     \/_/   \/_____/   \/_____/   \/_____/   \/_____/   \/_/ /_/   \/_/\/_/   \/_/  \/_/ 
                                                                                                                                                                                  `

	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	white := color.New(color.FgWhite, color.Bold)

	red.Println(banner)
	println()

	http_port := cfg.HttpPort

	yellow.Println("Configuration:")
	white.Printf("   Port:        %d\n", http_port)
	white.Printf("   JWT Expiry:  %s\n", time.Duration(cfg.Expiry).String())
	white.Printf("   Rate Limit:  Public=%d RPM, Protected=%d RPM\n", cfg.PublicRPM, cfg.ProtectedRPM)
	println()

	green.Println("Services:")
	green.Println("   ✓ PostgreSQL connected")
	green.Println("   ✓ Redis connected")
	green.Println("   ✓ JWT Auth middleware enabled")
	green.Println("   ✓ Rate limiting enabled")
	green.Println("   ✓ Logging middleware enabled")
	green.Println("   ✓ Recovery middleware enabled")
	println()

	magenta.Println("Endpoints:")
	white.Printf("   API Base:     http://localhost:%d/\n", http_port)
	white.Printf("   Swagger:      http://localhost:%d/swagger/\n", http_port)
	white.Printf("   Swagger JSON: http://localhost:%d/swagger/doc.json\n", http_port)
	println()

	yellow.Println("🚀 Server is running! Press CTRL+C to stop.")
}

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

	// connect to cache
	rdb, err := redis.RedisConnect(cfg)
	if err != nil {
		logger.Err(err).Str("component", "Redis").Msg("Redis could not connect to db.")
	}

	// connect to DB
	DB, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Str("component", "Postgres").Msg("Postgres could not connect to db.")
	}

	// prepare user controller
	hasher := hasher.NewBcryptHasher(0)
	userRepo := postgres.NewUserRepository(DB)
	userSVC := service.NewUserService(userRepo, hasher)
	userCtrl := handler.NewUserController(userSVC)

	// prepare auth controller
	prov := jwt.NewProvider(cfg.SecretKey, time.Duration(cfg.Expiry))
	authSVC := service.NewAuthService(userRepo, hasher, prov)
	authCtrl := handler.NewAuthController(authSVC)

	r := handler.NewRouter(rdb, userCtrl, authCtrl, cfg.SecretKey, time.Duration(cfg.Expiry), cfg.PublicRPM, cfg.ProtectedRPM, logger)

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:    strconv.FormatInt(int64(cfg.HttpPort), 10),
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

	printBanner(cfg)
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
