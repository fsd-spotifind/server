package database

import (
	"context"
	"fmt"
	"time"

	db "fsd-backend/prisma/db"
)

func (s *service) CreateUser(ctx context.Context, email, username, passwordHash string) (*db.UserModel, error) {
	return s.client.User.CreateOne(
		db.User.Username.Set(username),
		db.User.PasswordHash.Set(passwordHash),
		db.User.Email.Set(email),
		db.User.IsVerified.Set(false),
	).Exec(ctx)
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*db.UserModel, error) {
	return s.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)
}

func (s *service) CreateRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error {
	_, err := s.client.RefreshToken.CreateOne(
		db.RefreshToken.Token.Set(token),
		db.RefreshToken.User.Link(db.User.ID.Equals(userID)),
		db.RefreshToken.ExpiresAt.Set(expiresAt),
	).Exec(ctx)
	return err
}

func (s *service) GetUserFromRefreshToken(ctx context.Context, token string) (*db.UserModel, error) {
	refreshToken, err := s.client.RefreshToken.FindUnique(
		db.RefreshToken.Token.Equals(token),
	).With(
		db.RefreshToken.User.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Check if token is revoked
	revokedAt, ok := refreshToken.RevokedAt()
	if !ok {
		return refreshToken.User(), nil
	}

	// if RevokedAt is set, check if it's after creation (meaning it was revoked)
	if revokedAt.After(refreshToken.CreatedAt) {
		return nil, fmt.Errorf("refresh token revoked")
	}

	return refreshToken.User(), nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := s.client.RefreshToken.FindUnique(
		db.RefreshToken.Token.Equals(token),
	).Update(
		db.RefreshToken.RevokedAt.Set(time.Now().UTC()),
		db.RefreshToken.UpdatedAt.Set(time.Now().UTC()),
	).Exec(ctx)
	return err
}

// simple semi-flawed health check that works only if we have 1 user at least
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// find a single user to test connection
	_, err := s.client.User.FindFirst().Exec(ctx)
	if err != nil && err != db.ErrNotFound {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "Database connection is healthy"
	return stats
}

func (s *service) Close() error {
	return s.client.Prisma.Disconnect()
}
