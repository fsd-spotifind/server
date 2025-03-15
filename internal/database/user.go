package database

import (
	"context"
	"fmt"
	"time"

	db "fsd-backend/prisma/db"
)

func (s *service) GetUserByEmail(ctx context.Context, email string) (*db.UserModel, error) {
	return s.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)
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
