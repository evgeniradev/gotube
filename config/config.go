package config

import (
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
)

// Holds various application-wide configurations
type AppConfig struct {
	Session              *scs.SessionManager
	DB                   DatabaseOperations
	EnableCSRFProtection bool
	InfoLog              *log.Logger
	ErrorLog             *log.Logger
}

// Global instance of the app configuration
var App AppConfig = AppConfig{
	Session:              createSession(),
	EnableCSRFProtection: true,
	InfoLog:              log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	ErrorLog:             log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
}
