package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mrinalwahal/service/handler"
	"github.com/mrinalwahal/service/pkg/middleware"
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
	//With("release", "v1.0.0")

	//	Initialize the server.
	router := http.NewServeMux()

	handler := handler.NewHTTPHandler(&handler.HTTPHandlerConfig{
		// Prefix: "/v1",
		Logger: logger,
	})

	// Health check endpoint.
	// Returns OK if the server is running.
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// CRUD routes.
	router.HandleFunc("POST /v1/", handler.HandlerFunc(handler.Create))
	router.HandleFunc("GET /v1/{id}", handler.HandlerFunc(handler.Get))

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
	// s := server.NewHTTPServer(&server.NewHTTPServerConfig{
	// 	Port:      "8080",
	// 	Dialector: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"),
	// 	Logger:    logger,
	// })

	//	Start the server.
	server := http.Server{
		Addr:    ":8080",
		Handler: chain(router),
	}
	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()
}
