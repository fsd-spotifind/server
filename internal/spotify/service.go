package spotify

import (
	"context"
	"fmt"
	"fsd-backend/internal/models"
	"net/http"
	"strings"
)

type Service interface {
	GetTrackByID(ctx context.Context, token string, itemId string) (*models.Track, error)
	GetTracksByIDs(ctx context.Context, token string, ids []string) ([]models.Track, error)
	GetUserRecentlyPlayedTracks(ctx context.Context, token string) (*models.RecentlyPlayed, error)
}

type service struct {
	client *Client
}

func NewService(config *Config) Service {
	return &service{
		client: NewClient(config),
	}
}

func (s *service) GetTrackByID(ctx context.Context, token string, trackId string) (*models.Track, error) {
	req, err := s.client.request(ctx, http.MethodGet, "/tracks/"+trackId, token)
	if err != nil {
		fmt.Println("Error getting track by id", err)
		return nil, err
	}
	var track models.Track
	if err := s.client.doJSON(req, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func (s *service) GetTracksByIDs(ctx context.Context, token string, ids []string) ([]models.Track, error) {
	// Spotify API limit is 50 ids per request
	const batchSize = 50
	var allTracks []models.Track
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[start:end]
		idParam := strings.Join(batchIDs, ",")
		req, err := s.client.request(ctx, http.MethodGet, "/tracks?ids="+idParam, token)
		if err != nil {
			return nil, err
		}
		var resp struct {
			Tracks []models.Track `json:"tracks"`
		}
		if err := s.client.doJSON(req, &resp); err != nil {
			return nil, err
		}
		allTracks = append(allTracks, resp.Tracks...)
	}
	return allTracks, nil
}

func (s *service) GetUserRecentlyPlayedTracks(ctx context.Context, token string) (*models.RecentlyPlayed, error) {
	req, err := s.client.request(ctx, http.MethodGet, "/me/player/recently-played", token)
	if err != nil {
		fmt.Println("Error getting user recently played tracks", err)
		return nil, err
	}
	var recentlyPlayed models.RecentlyPlayed
	if err := s.client.doJSON(req, &recentlyPlayed); err != nil {
		return nil, err
	}
	return &recentlyPlayed, nil
}
