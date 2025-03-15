package database

import (
	"context"
	"log"
	"time"

	db "fsd-backend/prisma/db"
)

type Service interface {
	GetUserByEmail(ctx context.Context, email string) (*db.UserModel, error)
	CreateSotd(ctx context.Context, userID, trackID, note, mood string) (*db.SongOfTheDayModel, error)
	GetSotd(ctx context.Context, userId string, date time.Time) (*db.SongOfTheDayModel, error)
	UpdateSotd(ctx context.Context, sotdId, trackId, note, mood string) (*db.SongOfTheDayModel, error)
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
