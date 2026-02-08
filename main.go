package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/audunhov/member/auth"
	"github.com/audunhov/member/database/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func strPtr(s string) *string { return &s }

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	fmt.Println("URL:",dbUrl)
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	queries := db.New(conn)

	as := auth.NewAuthService(conn, queries)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := as.RegisterUser(r.Context(), email, password)
		if err != nil {
			http.Error(w, "Could not log in: "+err.Error(), http.StatusBadRequest)
			return
		}
		user, err := as.Login(r.Context(), email, password)
		if err != nil {
			http.Error(w, "This should never happen", 500)
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := as.Login(r.Context(), email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.String())
		w.Write([]byte("Hello world"))
	})
	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", mux)
}
