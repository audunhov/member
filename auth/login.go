package auth

import (
	"context"
	"fmt"

	"github.com/audunhov/member/database/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewAuthService(pool *pgxpool.Pool, q *db.Queries) AuthService {
	return AuthService{
		pool: pool,
		q:    q,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string) error {
	hashedPass, err := hashPassword(password)
	if err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := s.q.WithTx(tx)

	member, err := q.CreateMember(ctx, db.CreateMemberParams{
		Email: email,
		Data:  pgtype.Map{},
	})
	if err != nil {
		return fmt.Errorf("failed creating member: %w", err)
	}
	err = q.CreateLocalAuth(ctx, db.CreateLocalAuthParams{
		MemberID:     member.ID,
		PasswordHash: hashedPass,
	})
	if err != nil {
		return fmt.Errorf("Failed creating auth: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*db.Member, error) {
	result, err := s.q.GetMemberWithPassword(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	if !checkPasswordHash(password, result.PasswordHash) {
		return nil, fmt.Errorf("Invalid email or password")
	}

	member := db.Member{
		ID:              result.ID,
		Email:           result.Email,
		Data:            result.Data,
		EmailVerifiedAt: result.EmailVerifiedAt,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
	}

	return &member, nil
}
