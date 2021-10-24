package config

import (
	"github.com/kkyr/fig"
	"log"
	"os"
)

// Config structure for all settings of application
type Config struct {
	Service struct {
		ServerAddress    string `fig:"serverAddress" default:"localhost"`    // address of server
		ServerPort       int    `fig:"serverPort" default:"4112"`            // port of server
		TypeOfService    string `fig:"typeOfService" default:"stdlib"`       // type of service: (stdlib, gin, gorilla, fasthttp)
		TypeOfRepository string `fig:"typeOfRepository" default:"in-memory"` // type of repository: (in-memory, PostrgreSQL, MySQL)
	} `fig:"service"`
	ErrorLogger *log.Logger // logger for use, don't load from configuration file
}

// Init function for initialize Config structure
func Init() (*Config, error) {
	var cfg = Config{}
	err := fig.Load(&cfg, fig.Dirs("../", "./", "./..."), fig.File("config.yaml"))
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
		return nil, err
	}

	cfg.ErrorLogger = NewBuiltinLogger().logger

	return &cfg, err
}

// BuiltinLogger custom logger
type BuiltinLogger struct {
	logger *log.Logger
}

// NewBuiltinLogger function custom logger initialize
func NewBuiltinLogger() *BuiltinLogger {
	return &BuiltinLogger{logger: log.New(os.Stdout, "", 5)}
}

// Debug method for print debug
func (l *BuiltinLogger) Debug(args ...interface{}) {
	l.logger.Println(args...)
}

// Debugf method for print formatted debug
func (l *BuiltinLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
