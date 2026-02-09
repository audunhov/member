package api

import (
	"github.com/audunhov/member/database/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/alexedwards/scs/v2"
)

type Server struct {
	DB *db.Queries
	Pool *pgxpool.Pool
	Session *scs.SessionManager
}

func NewServer(pool *pgxpool.Pool, session *scs.SessionManager) *Server {
	return &Server{
		DB:      db.New(pool),
		Pool:    pool,
		Session: session,
	}
}
