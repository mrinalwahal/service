package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mrinalwahal/service/pkg/middleware"
	"github.com/mrinalwahal/service/router/http/router"
	"gorm.io/driver/postgres"
)

func main() {

	err := godotenv.Load(".env.example")
	if err != nil {
		log.Println("Error loading .env.development file")
	}

	//	Setup the logger.
	level := slog.LevelInfo
	DEBUG, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		panic(err)
	}
	if DEBUG {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		Level: level,
	}))
	logger = logger.
		With("service", "todo").
		With("environment", os.Getenv("ENV"))

	//	Initialize the router.
	router := router.NewHTTPRouter(&router.HTTPRouterConfig{
		Dialector: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"),
		Logger:    logger,
	})

	// Prepare the middleware chain.
	// The order of the middlewares is important.
	// Recommended order: Request ID -> RateLimit -> CORS -> Logging -> Recover -> Auth -> Cache -> Compression
	chain := middleware.Chain(
		middleware.RequestID,
		// TODO: middleware.RateLimit,
		middleware.CORS,
		middleware.Recover(logger),
		middleware.Logging(logger),
	)

	// Prepare the base router.
	baseRouter := http.NewServeMux()
	baseRouter.Handle("/todo/", http.StripPrefix("/todo", router))

	//	Configure and start the server.
	server := http.Server{
		Addr:     ":8080",
		Handler:  chain(baseRouter),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()
}
