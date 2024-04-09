package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/pkg/middleware"
	"github.com/mrinalwahal/service/service"
	"github.com/mrinalwahal/service/transport/http/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
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
		With("service", "record").
		With("environment", os.Getenv("ENV"))

	//	Setup the gorm logger.
	handler := logger.With("layer", "database").Handler()
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(handler), // since v1.3.0
		slogGorm.WithTraceAll(),       // trace all messages
	)

	// Open a database connection.
	conn, err := gorm.Open(postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		panic(err)
	}

	// Configure connection pooling.
	//
	// Link: https://gorm.io/docs/generic_interface.html#Connection-Pool
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	// Connect the database layer.
	db := db.NewDB(&db.Config{
		DB: conn,
	})

	// GORM provides Prometheus plugin to collect DBStats or user-defined metrics
	// https://gorm.io/docs/prometheus.html
	// https://github.com/go-gorm/prometheus
	//
	// db.Use(prometheus.New(prometheus.Config{
	// 	DBName:          "db1",                       // use `DBName` as metrics label
	// 	RefreshInterval: 15,                          // Refresh metrics interval (default 15 seconds)
	// 	PushAddr:        "prometheus pusher address", // push metrics if `PushAddr` configured
	// 	StartServer:     true,                        // start http server to expose metrics
	// 	HTTPServerPort:  8080,                        // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
	// 	MetricsCollector: []prometheus.MetricsCollector{
	// 		&prometheus.MySQL{
	// 			VariableNames: []string{"Threads_running"},
	// 		},
	// 	}, // user defined metrics
	// }))

	// Get the service layer.
	service := service.NewService(&service.Config{
		DB:     db,
		Logger: logger,
	})

	//	Initialize the router.
	router := router.NewHTTPRouter(&router.HTTPRouterConfig{
		Service: service,
		Logger:  logger,
	})

	// Prepare the middleware chain.
	// The order of the middlewares is important.
	// Recommended order: Request ID -> RateLimit -> CORS -> Logging -> Recover -> Auth -> Cache -> Compression
	chain := middleware.Chain(
		middleware.RequestID,
		// TODO: middleware.RateLimit,
		middleware.CORS,
		middleware.Recover(logger.With("server", "http")),
		middleware.Logging(logger.With("server", "http")),
	)

	// Prepare the base router.
	baseRouter := http.NewServeMux()
	baseRouter.Handle("/records/", http.StripPrefix("/records", router))

	//	Configure and start the server.
	server := http.Server{
		Addr:     ":8080",
		Handler:  chain(baseRouter),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()

	// Close the database connection.
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}
