package server

import (
	"context"
	"fmt"

	"fsd-backend/internal/utils"
	db "fsd-backend/prisma/db"
)

func (s *Server) GenerateUserStatistic(ctx context.Context, userID string, accessToken string, period db.StatisticPeriod) error {
	since := utils.PeriodToStartTime(period)
	recentlyPlayed, err := s.spotify.GetUserRecentlyPlayedTracks(ctx, accessToken, &since)
	if err != nil {
		return fmt.Errorf("failed to get recently played tracks: %w", err)
	}
	computedStats := utils.ComputeUserStatistics(recentlyPlayed)
	_, err = s.db.CreateUserStatistic(ctx, userID, period,
		computedStats.TotalTracks,
		computedStats.TotalDuration,
		computedStats.UniqueArtists,
		computedStats.Vibe,
		computedStats.TopArtistIDs,
		computedStats.TopTrackIDs,
		computedStats.TopAlbumIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to create user statistic: %w", err)
	}
	return nil
}

func (s *Server) RefreshAccessTokenIfNeeded(ctx context.Context, user *db.AccountModel) (string, error) {

	refreshToken, _ := user.RefreshToken()
	result, err := s.spotify.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken := result.AccessToken
	err = s.db.UpdateAccessToken(ctx, user.ID, result.AccessToken, result.RefreshToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
