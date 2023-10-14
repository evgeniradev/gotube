package config

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// Creates and configures a new session manager
func createSession() *scs.SessionManager {
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	return session
}
