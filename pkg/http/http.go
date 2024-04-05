package http

import (
	"log/slog"
	"net/http"
)

type HTTPServerConfig struct {

	//	Port on which the server will listen on.
	//	Example: ":8080"
	//
	//	Note: If the port is not prefixed with `:`, it will be prefixed automatically.
	//	Example: "8080" will be converted to ":8080"
	//
	//	This field is required.
	Port string

	//	Host on which the server will listen on.
	//	Example: "localhost"
	//	Default: "localhost"
	//
	//	This field is optional.
	Host string

	//	Logger is the `log/slog` instance that will be used to log messages.
	//	Default: `slog.DefaultLogger`
	//
	//	This field is optional.
	Logger *slog.Logger
}

type HTTPServer struct {
	*http.Server
	log *slog.Logger
}

// NewHTTPServer creates a new instance of `HTTPServer`.
func NewHTTPServer(config *HTTPServerConfig) *HTTPServer {
	if config.Port == "" {
		panic("Port is required")
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return &HTTPServer{
		Server: &http.Server{
			Addr:     config.Host + config.Port,
			Handler:  http.NewServeMux(),
			ErrorLog: slog.NewLogLogger(config.Logger.Handler(), slog.LevelError),
		},
		log: config.Logger,
	}
}

// Serve starts the HTTP server.
func (s *HTTPServer) ListenAndServe() {
	s.log.Info("HTTP server started on " + s.Addr)
	s.Server.ListenAndServe()
}
