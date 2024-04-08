package test

import (
	"log/slog"
	"testing"

	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"
)

// environment contains all the configuration that is required by our tests.
type environment struct {

	// Database connection.
	db *gorm.DB

	// Logger instance.
	log *slog.Logger

	// Service layer.
	service service.Service
}

// initialize configures a suitable and reliable environment for the tests.
func initialize(t *testing.T) *environment {
	return nil
}
