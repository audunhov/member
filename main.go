package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/audunhov/member/internal/api"
	"github.com/jackc/pgx/v5/pgxpool"
)


func setupSessions(pool *pgxpool.Pool) *scs.SessionManager {
	sm := scs.New()

	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 24 * time.Hour
	sm.IdleTimeout = 1 * time.Hour

	sm.Cookie.Name = "session_id"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Secure = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Path = "/"

	return sm
}

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		slog.Error("Unable to connect to database", "err", err)
		return
	}
	defer conn.Close()

	sm := setupSessions(conn)

	srv := api.NewServer(conn, sm)

	server := http.Server{
		Addr:    ":8080",
		Handler: srv.Routes(),
	}

	fmt.Println("Starting server at :8080")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server crashed")
	}
}
