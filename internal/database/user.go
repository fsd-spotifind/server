package database

import (
	"context"
	"fmt"
	"time"

	db "fsd-backend/prisma/db"
)

func (s *service) GetUserById(ctx context.Context, userId string) (*db.UserModel, error) {
	return s.client.User.FindUnique(
		db.User.ID.Equals(userId),
	).Exec(ctx)
}

func (s *service) GetUserAccountByUserId(ctx context.Context, userId string) (*db.AccountModel, error) {
	return s.client.Account.FindFirst(
		db.Account.UserID.Equals(userId),
	).Exec(ctx)
}

func (s *service) GetAllUserAccounts(ctx context.Context) ([]db.AccountModel, error) {
	accounts, err := s.client.Account.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *service) CreateUserStatistic(ctx context.Context, userId string, period db.StatisticPeriod, totalTracks, totalDuration, uniqueArtists int, vibe string, topArtistsIds, topTracksIds, topAlbumsIds []string) (*db.UserStatisticModel, error) {
	return s.client.UserStatistic.CreateOne(
		db.UserStatistic.Period.Set(period),
		db.UserStatistic.User.Link(db.User.ID.Equals(userId)),
		db.UserStatistic.TotalTracks.Set(totalTracks),
		db.UserStatistic.TotalDuration.Set(totalDuration),
		db.UserStatistic.UniqueArtists.Set(uniqueArtists),
		db.UserStatistic.Vibe.Set(vibe),
		db.UserStatistic.TopArtistsIds.Set(topArtistsIds),
		db.UserStatistic.TopTracksIds.Set(topTracksIds),
		db.UserStatistic.TopAlbumsIds.Set(topAlbumsIds),
	).Exec(ctx)
}

func (s *service) GetUserStatisticByPeriod(ctx context.Context, userId string, period db.StatisticPeriod) ([]db.UserStatisticModel, error) {
	return s.client.UserStatistic.FindMany(
		db.UserStatistic.UserID.Equals(userId),
		db.UserStatistic.Period.Equals(period),
	).Exec(ctx)
}

func (s *service) UpdateAccessToken(ctx context.Context, userId string, accessToken string) error {
	_, err := s.client.Account.FindUnique(
		db.Account.ID.Equals(userId),
	).Update(
		db.Account.AccessToken.Set(accessToken),
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
