package main

import (
	"log/slog"

	"github.com/grafana/loki-client-go/loki"
	"github.com/mrinalwahal/service/server"
	slogloki "github.com/samber/slog-loki/v3"
	"gorm.io/driver/postgres"
)

func main() {

	//	Setup the loki client to use loki as log sink.
	config, _ := loki.NewDefaultConfig("http://localhost:3100/loki/api/v1/push")
	//	config.TenantID = "xyz"
	client, _ := loki.New(config)
	defer client.Stop()

	//	Setup the logger.
	logger := slog.New(slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler())
	logger = logger.
		With("service", "todo")
		//With("environment", "dev").
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
