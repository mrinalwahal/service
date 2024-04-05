package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/grafana/loki-client-go/loki"
	"github.com/mrinalwahal/service/server"
	slogloki "github.com/samber/slog-loki/v3"
	"gorm.io/driver/postgres"
)

func main() {

	//	Setup the loki client to use loki as log sink.
	config, _ := loki.NewDefaultConfig(fmt.Sprintf("%s/loki/api/v1/push", os.Getenv("LOKI_HOST")))
	//	config.TenantID = "xyz"
	client, _ := loki.New(config)
	defer client.Stop()

	//	Setup the logger.
	level := slog.LevelInfo
	DEBUG, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		panic(err)
	}
	if DEBUG {
		level = slog.LevelDebug
	}

	logger := slog.New(slogloki.Option{Level: level, Client: client}.NewLokiHandler())
	logger = logger.
		With("service", "todo").
		With("environment", os.Getenv("ENV"))
	//With("release", "v1.0.0")

	//	Initialize the server.
	s := server.NewHTTPServer(&server.NewHTTPServerConfig{
		Port:      "8080",
		Dialector: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"),
		Logger:    logger,
	})

	//	Start the server.
	s.Serve()
}
