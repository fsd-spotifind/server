package database

import (
	"context"
	"fmt"
	"log"
	"time"

	db "fsd-backend/prisma/db"
)

type Service interface {
	CreateUser(ctx context.Context, email, username, passwordHash string) (*db.UserModel, error)
	Health() map[string]string
	Close() error
}

type service struct {
	client *db.PrismaClient
}

var dbInstance *service

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		client: client,
	}
	return dbInstance
}

func (s *service) CreateUser(ctx context.Context, email, username, passwordHash string) (*db.UserModel, error) {
	return s.client.User.CreateOne(
		db.User.Username.Set(username),
		db.User.PasswordHash.Set(passwordHash),
		db.User.Email.Set(email),
		db.User.IsVerified.Set(false),
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
