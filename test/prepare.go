package test

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/moritiza/url-shortener/config"
	"github.com/sirupsen/logrus"
)

// Prepare for testing application
func Prepare() config.Config {
	var cfg config.Config

	// Create new logrus
	cfg.Logger = logrus.New()
	cfg.Logger.SetReportCaller(true)
	cfg.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf(" *** Function: %s() *** Message:", f.Function),
				fmt.Sprintf(" *** File: %s *** Line: %d", filepath.Base(f.File), f.Line)
		},
	})

	// Create new validator
	cfg.Validator = config.Validator()

	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		cfg.Logger.Fatalf("Could not load environments!\n%v", err)
	}

	// Connect to database
	cfg.Database = config.ConnectTestDatabase(cfg)
	return cfg
}
