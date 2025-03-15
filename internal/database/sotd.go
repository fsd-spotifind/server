package database

import (
	"context"
	db "fsd-backend/prisma/db"
	"time"
)

const (
	YYYYMMDD = "2006-01-02"
)

func (s *service) CreateSotd(ctx context.Context, userId, trackId, note, mood string) (*db.SongOfTheDayModel, error) {
	return s.client.SongOfTheDay.CreateOne(
		db.SongOfTheDay.TrackID.Set(trackId),
		db.SongOfTheDay.SetAt.Set(time.Now().Format(YYYYMMDD)),
		db.SongOfTheDay.User.Link(db.User.ID.Equals(userId)),
		db.SongOfTheDay.Note.Set(note),
		db.SongOfTheDay.Mood.Set(mood),
	).Exec(ctx)
}

func (s *service) GetSotd(ctx context.Context, userId string, date time.Time) (*db.SongOfTheDayModel, error) {

	dateString := date.Format(YYYYMMDD)

	sotd, err := s.client.SongOfTheDay.FindFirst(
		db.SongOfTheDay.SetAt.Equals(dateString),
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
