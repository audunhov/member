package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type Middleware func(next http.Handler) http.Handler

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.Session.Exists(r.Context(), "userID") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				slog.Error("Panic recovered", "err", err, "stack", debug.Stack())
				errorResponse(w, http.StatusInternalServerError, "Unexpected error occured")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.RemoteAddr, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Routes() http.Handler {
	var handler http.Handler
	mux := http.NewServeMux()

	api := http.NewServeMux()

	v1 := http.NewServeMux()
	v1.HandleFunc("POST /login", s.handleLogin)
	v1.HandleFunc("POST /register", s.handleRegister)
	v1.Handle("GET /auth/me", s.AuthMiddleware(http.HandlerFunc(s.handleMe)))

	api.Handle("/v1/", http.StripPrefix("/v1", v1))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	mux.Handle("/api/", http.StripPrefix("/api", api))

	middleware := []Middleware{
		s.Session.LoadAndSave,
		loggerMiddleware,
		recoverMiddleware,
	}

	handler = mux
	for _, mw := range middleware {
		handler = mw(handler)
	}

	fmt.Println(handler)

	return handler
}
