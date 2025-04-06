package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"fsd-backend/internal/models"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service interface {
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetTrackByID(ctx context.Context, token string, itemId string) (*models.Track, error)
	GetTracksByIDs(ctx context.Context, token string, ids []string) ([]models.Track, error)
	GetAlbumsByIDs(ctx context.Context, token string, ids []string) ([]models.Album, error)
	GetArtistsByIDs(ctx context.Context, token string, ids []string) ([]models.Artist, error)
	GetUserRecentlyPlayedTracks(ctx context.Context, token string, since *time.Time) (*models.RecentlyPlayed, error)
}

type service struct {
	client *Client
}

func NewService(config *Config) Service {
	return &service{
		client: NewClient(config),
	}
}

func (s *service) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.client.config.ClientID, s.client.config.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("refresh failed: %s", string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
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

func (s *service) GetAlbumsByIDs(ctx context.Context, token string, ids []string) ([]models.Album, error) {
	// Spotify API limit is 20 ids per request
	const batchSize = 20
	var allAlbums []models.Album
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[start:end]
		idParam := strings.Join(batchIDs, ",")
		req, err := s.client.request(ctx, http.MethodGet, "/albums?ids="+idParam, token)
		if err != nil {
			return nil, err
		}
		var resp struct {
			Albums []models.Album `json:"albums"`
		}
		if err := s.client.doJSON(req, &resp); err != nil {
			return nil, err
		}
		allAlbums = append(allAlbums, resp.Albums...)
	}
	return allAlbums, nil
}

func (s *service) GetArtistsByIDs(ctx context.Context, token string, ids []string) ([]models.Artist, error) {
	// Spotify API limit is 50 ids per request
	const batchSize = 50
	var allArtists []models.Artist
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[start:end]
		idParam := strings.Join(batchIDs, ",")
		req, err := s.client.request(ctx, http.MethodGet, "/artists?ids="+idParam, token)
		if err != nil {
			return nil, err
		}
		var resp struct {
			Artists []models.Artist `json:"artists"`
		}
		if err := s.client.doJSON(req, &resp); err != nil {
			return nil, err
		}
		allArtists = append(allArtists, resp.Artists...)
	}
	return allArtists, nil
}

func (s *service) GetUserRecentlyPlayedTracks(ctx context.Context, token string, since *time.Time) (*models.RecentlyPlayed, error) {
	endpoint := "/me/player/recently-played"
	if since != nil {
		after := since.UnixMilli()
		endpoint = fmt.Sprintf("%s?after=%d", endpoint, after)
	}
	req, err := s.client.request(ctx, http.MethodGet, endpoint, token)
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
