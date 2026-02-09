package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/audunhov/member/auth"
	"github.com/audunhov/member/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSON[loginRequest](r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	authData, err := s.DB.GetMemberWithPassword(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errorResponse(w, http.StatusUnauthorized, "Wrong email or password")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	ok := auth.CheckPasswordHash(req.Password, authData.PasswordHash)
	if ok != true {
		errorResponse(w, http.StatusUnauthorized, "Wrong email or password")
		return
	}

	err = s.Session.RenewToken(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not create session")
		return
	}

	s.Session.Put(r.Context(), "userID", authData.ID.String())

	jsonResponse(w, http.StatusOK, loginResponse{
		UserID: authData.ID.String(),
		Email:  authData.Email,
	})
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type registerResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSON[registerRequest](r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	hashedPass, err := auth.HashPassword(req.Password)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	tx, err := s.Pool.Begin(r.Context())
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Bad request")
		return

	}
	defer tx.Rollback(r.Context())

	q := s.DB.WithTx(tx)

	member, err := q.CreateMember(r.Context(), db.CreateMemberParams{
		Email: req.Email,
		Data:  pgtype.Map{},
	})
	if err != nil {
		return
	}

	err = q.CreateLocalAuth(r.Context(), db.CreateLocalAuthParams{
		MemberID:     member.ID,
		PasswordHash: hashedPass,
	})
	if err != nil {
		fmt.Errorf("Failed creating auth: %w", err)
		return
	}
	err = s.Session.RenewToken(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not create session")
		return
	}

	s.Session.Put(r.Context(), "userID", member.ID.String())


	tx.Commit(r.Context())

	jsonResponse(w, http.StatusCreated, registerResponse{
		UserID: member.ID.String(),
		Email:  member.Email,
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	id, ok := s.Session.Get(r.Context(), "userID").(string)
	if !ok {
		errorResponse(w, http.StatusInternalServerError, "Invalid session, should be blocked in middleware")
		return
	}
	user, err := s.DB.GetMemberById(r.Context(), uuid.MustParse(id))
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}
	jsonResponse(w, http.StatusOK, user)
}
