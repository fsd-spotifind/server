package database

import (
	"context"
	db "fsd-backend/prisma/db"
	"time"
)

const (
	YYYYMMDD = "2006-01-02"
)

type PaginatedSotdResponse struct {
	SotdEntries []db.SongOfTheDayModel `json:"sotdEntries"`
	HasMore     bool                   `json:"hasMore"`
}

func (s *service) CreateSotd(ctx context.Context, userId, trackId, note, mood string) (*db.SongOfTheDayModel, error) {
	return s.client.SongOfTheDay.CreateOne(
		db.SongOfTheDay.TrackID.Set(trackId),
		db.SongOfTheDay.SetAt.Set(time.Now().Format(YYYYMMDD)),
		db.SongOfTheDay.User.Link(db.User.ID.Equals(userId)),
		db.SongOfTheDay.Note.Set(note),
		db.SongOfTheDay.Mood.Set(mood),
	).Exec(ctx)
}

func (s *service) GetSotdByDate(ctx context.Context, userId string, date time.Time) (*db.SongOfTheDayModel, error) {

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

func (s *service) GetAllSotd(ctx context.Context, userId string, limit, offset int) (*PaginatedSotdResponse, error) {
	sotdEntries, err := s.client.SongOfTheDay.FindMany(
		db.SongOfTheDay.UserID.Equals(userId),
	).
		Take(limit + 1).
		Skip(offset).
		OrderBy(
			db.SongOfTheDay.SetAt.Order(db.DESC),
		).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	hasMore := len(sotdEntries) > limit

	if hasMore {
		sotdEntries = sotdEntries[:limit]
	}

	return &PaginatedSotdResponse{
		SotdEntries: sotdEntries,
		HasMore:     hasMore,
	}, nil
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
