package database

import (
	"context"
	"errors"
	db "fsd-backend/prisma/db"
	"time"
)

func (s *service) CreateSotd(ctx context.Context, userId, trackId, note, mood string) (*db.SongOfTheDayModel, error) {
	currentDate := time.Now().UTC()
	sotd, err := s.GetSotd(ctx, userId, currentDate)
	if err != nil {
		return nil, err
	}
	if sotd != nil {
		return nil, errors.New("SOTD already exists, please update it instead of creating a new one")
	}
	return s.client.SongOfTheDay.CreateOne(
		db.SongOfTheDay.TrackID.Set(trackId),
		db.SongOfTheDay.User.Link(db.User.ID.Equals(userId)),
		db.SongOfTheDay.Note.Set(note),
		db.SongOfTheDay.Mood.Set(mood),
		db.SongOfTheDay.SetAt.Set(time.Now()),
	).Exec(ctx)
}

func (s *service) GetSotd(ctx context.Context, userId string, date time.Time) (*db.SongOfTheDayModel, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	sotd, err := s.client.SongOfTheDay.FindFirst(
		db.SongOfTheDay.SetAt.Gte(startOfDay),
		db.SongOfTheDay.SetAt.Lte(endOfDay),
		db.SongOfTheDay.UserID.Equals(userId),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return sotd, nil
}

func (s *service) UpdateSotd(ctx context.Context, sotdId, trackId, note, mood string) (*db.SongOfTheDayModel, error) {
	return s.client.SongOfTheDay.FindUnique(
		db.SongOfTheDay.ID.Equals(sotdId),
	).Update(
		db.SongOfTheDay.TrackID.Set(trackId),
		db.SongOfTheDay.Note.Set(note),
		db.SongOfTheDay.Mood.Set(mood),
	).Exec(ctx)
}
