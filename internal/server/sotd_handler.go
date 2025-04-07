package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fsd-backend/internal/models"
)

func (s *Server) createSotdHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")

	type parameters struct {
		TrackID string `json:"track_id"`
		Note    string `json:"note"`
		Mood    string `json:"mood"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	sotd, err := s.db.CreateSotd(r.Context(), UserID, params.TrackID, params.Note, params.Mood)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create Sotd", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, sotd)
}

func (s *Server) getSotdBydateHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	Date := r.PathValue("date")

	UnixTime, err := strconv.ParseInt(Date, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Unix timestamp format", http.StatusBadRequest)
		return
	}

	ParsedDate := time.Unix(UnixTime, 0).UTC()

	sotd, err := s.db.GetSotdByDate(r.Context(), UserID, ParsedDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Sotd", err)
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

	track, err := s.spotify.GetTrackByID(r.Context(), accessToken, sotd.TrackID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Track", err)
		return
	}

	note, _ := sotd.Note()
	mood, _ := sotd.Mood()
	formattedSotd := models.SotdEntry{
		ID:        sotd.ID,
		UserID:    sotd.UserID,
		TrackID:   sotd.TrackID,
		Note:      note,
		Mood:      mood,
		SetAt:     sotd.SetAt,
		CreatedAt: sotd.CreatedAt.String(),
		UpdatedAt: sotd.UpdatedAt.String(),
		Track:     *track,
	}

	respondWithJSON(w, http.StatusOK, formattedSotd)
}

func (s *Server) getSotdsHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	sotds, err := s.db.GetAllSotd(r.Context(), UserID, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch Song of the Day entries", err)
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

	trackIDs := make([]string, len(sotds.SotdEntries))
	for i, entry := range sotds.SotdEntries {
		trackIDs[i] = entry.TrackID
	}

	tracks, err := s.spotify.GetTracksByIDs(r.Context(), accessToken, trackIDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch tracks", err)
		return
	}

	trackMap := make(map[string]models.Track, len(tracks))
	for _, track := range tracks {
		trackMap[track.ID] = track
	}

	grouped := make(map[string][]models.SotdEntry)
	for _, entry := range sotds.SotdEntries {
		note, _ := entry.Note()
		mood, _ := entry.Mood()
		grouped[entry.SetAt] = append(grouped[entry.SetAt], models.SotdEntry{
			ID:        entry.ID,
			UserID:    entry.UserID,
			TrackID:   entry.TrackID,
			Note:      note,
			Mood:      mood,
			SetAt:     entry.SetAt,
			CreatedAt: entry.CreatedAt.Format(time.RFC3339),
			UpdatedAt: entry.UpdatedAt.Format(time.RFC3339),
			Track:     trackMap[entry.TrackID],
		})
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"hasMore": sotds.HasMore,
		"entries": grouped,
	})
}

func (s *Server) updateSotdHandler(w http.ResponseWriter, r *http.Request) {
	SotdID := r.PathValue("sotdId")

	type parameters struct {
		TrackID string `json:"track_id"`
		Note    string `json:"note"`
		Mood    string `json:"mood"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	sotd, err := s.db.UpdateSotd(r.Context(), SotdID, params.TrackID, params.Note, params.Mood)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update Sotd", err)
		return
	}

	respondWithJSON(w, http.StatusOK, sotd)
}

func (s *Server) getRecommendedSotdsHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")

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

	recentlyPlayed, err := s.spotify.GetUserRecentlyPlayedTracks(r.Context(), accessToken, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Recently Played Tracks", err)
		return
	}

	respondWithJSON(w, http.StatusOK, recentlyPlayed)
}
