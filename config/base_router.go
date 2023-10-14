package config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

func BaseRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(SessionLoadMiddleware)
	router.Use(loadNoSurfMiddleware)
	router.Use(logRequestMiddleware)
	return router
}

// Loads and saves user sessions
func SessionLoadMiddleware(next http.Handler) http.Handler {
	return App.Session.LoadAndSave(next)
}

// Sets up CSRF protection using nosurf
func loadNoSurfMiddleware(next http.Handler) http.Handler {
	if !App.EnableCSRFProtection {
		return next
	}

	// Create a new nosurf handler
	csrfHandler := nosurf.New(next)

	// Configure the base CSRF cookie settings
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// Logs incoming HTTP requests
func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details including remote address, method, and URL
		App.InfoLog.Printf("%s - %s %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.Form, r.URL.RequestURI())

		// Pass the request to the next middleware or handler
		next.ServeHTTP(w, r)
	})
}
