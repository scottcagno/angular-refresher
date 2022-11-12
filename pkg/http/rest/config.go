package rest

import (
	"net/http"

	"github.com/cagnosolutions/webapp/pkg/webapp/logging"
	"github.com/cagnosolutions/webapp/pkg/webapp/middleware"
)

type Config struct {
	StaticHandler http.Handler `json:"-"`
	ErrHandler    http.Handler `json:"-"`
	MetricsOn     bool         `json:"metrics_on"`
	LoggingLevel  int          `json:"logging_level"`
}

var defaultConfig = &Config{
	StaticHandler: middleware.HandleStatic("/static/", "web/static/"),
	ErrHandler:    middleware.HandleErrors(),
	MetricsOn:     false,
	LoggingLevel:  logging.LevelInfo,
}

func checkConfig(conf *Config) {
	if conf == nil {
		conf = defaultConfig
	}
}
