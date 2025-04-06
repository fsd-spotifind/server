package server

import (
	"encoding/json"
	"fmt"
	"fsd-backend/internal/models"
	db "fsd-backend/prisma/db"
	"net/http"
	"strings"
)

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	user, err := s.db.GetUserById(r.Context(), UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user profile details", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (s *Server) createUserStatisticHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	type parameters struct {
		Period string `json:"period"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	period := db.StatisticPeriod(strings.ToUpper(params.Period))

	user, err := s.db.GetUserAccountByUserId(r.Context(), UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get User", err)
		return
	}
	accessToken, err := s.RefreshAccessTokenIfNeeded(r.Context(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get access token", err)
		return
	}

	if err := s.GenerateUserStatistic(r.Context(), UserID, accessToken, period); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate user statistic", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "User statistic created successfully",
	})
}

func (s *Server) getUserStatisticByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	Period := r.PathValue("period")

	userStatistic, err := s.db.GetUserStatisticByPeriod(r.Context(), UserID, db.StatisticPeriod(strings.ToUpper(Period)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user statistics", err)
		return
	}

	user, err := s.db.GetUserAccountByUserId(r.Context(), UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get User", err)
		return
	}
	accessToken, err := s.RefreshAccessTokenIfNeeded(r.Context(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get access token", err)
		return
	}

	var tracks []models.Track
	var artists []models.Artist
	var albums []models.Album

	allTrackIDs := make([]string, 0)
	allArtistIDs := make([]string, 0)
	allAlbumIDs := make([]string, 0)

	for _, stat := range userStatistic {
		allTrackIDs = append(allTrackIDs, stat.TopTracksIds...)
		allArtistIDs = append(allArtistIDs, stat.TopArtistsIds...)
		allAlbumIDs = append(allAlbumIDs, stat.TopAlbumsIds...)
	}

	if len(allTrackIDs) > 0 {
		fmt.Println("allTrackIDs", allTrackIDs)
		tracks, err = s.spotify.GetTracksByIDs(r.Context(), accessToken, allTrackIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch tracks", err)
			return
		}
	}

	if len(allArtistIDs) > 0 {
		artists, err = s.spotify.GetArtistsByIDs(r.Context(), accessToken, allArtistIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch artists", err)
			return
		}
	}

	if len(allAlbumIDs) > 0 {
		albums, err = s.spotify.GetAlbumsByIDs(r.Context(), accessToken, allAlbumIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch albums", err)
			return
		}
	}

	// Create maps for quick lookup
	trackMap := make(map[string]models.Track, len(tracks))
	for _, track := range tracks {
		trackMap[track.ID] = track
	}

	artistMap := make(map[string]models.Artist, len(artists))
	for _, artist := range artists {
		artistMap[artist.ID] = artist
	}

	albumMap := make(map[string]models.Album, len(albums))
	for _, album := range albums {
		albumMap[album.ID] = album
	}

	enhancedStats := make([]models.UserStatistic, len(userStatistic))
	for i, stat := range userStatistic {
		topTracks := make([]models.Track, 0, len(stat.TopTracksIds))
		for _, trackID := range stat.TopTracksIds {
			if track, ok := trackMap[trackID]; ok {
				topTracks = append(topTracks, track)
			}
		}
		topArtists := make([]models.Artist, 0, len(stat.TopArtistsIds))
		for _, artistID := range stat.TopArtistsIds {
			if artist, ok := artistMap[artistID]; ok {
				topArtists = append(topArtists, artist)
			}
		}
		topAlbums := make([]models.Album, 0, len(stat.TopAlbumsIds))
		for _, albumID := range stat.TopAlbumsIds {
			if album, ok := albumMap[albumID]; ok {
				topAlbums = append(topAlbums, album)
			}
		}
		vibe, _ := stat.Vibe()
		enhancedStats[i] = models.UserStatistic{
			ID:            stat.ID,
			UserID:        stat.UserID,
			Period:        stat.Period,
			TotalTracks:   stat.TotalTracks,
			TotalDuration: stat.TotalDuration,
			UniqueArtists: stat.UniqueArtists,
			Vibe:          vibe,
			TopTracks:     topTracks,
			TopArtists:    topArtists,
			TopAlbums:     topAlbums,
			CreatedAt:     stat.CreatedAt.String(),
			UpdatedAt:     stat.UpdatedAt.String(),
		}
	}
	respondWithJSON(w, http.StatusOK, enhancedStats)
}
