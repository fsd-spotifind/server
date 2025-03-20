package server

import (
	"encoding/json"
	"net/http"
	"strings"

	db "fsd-backend/prisma/db"
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
		Period        db.StatisticPeriod `json:"period"`
		TotalTracks   int                `json:"total_tracks"`
		TotalDuration int                `json:"total_duration"`
		UniqueArtists int                `json:"unique_artists"`
		Vibe          string             `json:"vibe"`
		TopArtistsIds []string           `json:"top_artists_ids"`
		TopTracksIds  []string           `json:"top_tracks_ids"`
		TopAlbumsIds  []string           `json:"top_albums_ids"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userStatistic, err := s.db.CreateUserStatistic(r.Context(), UserID, params.Period, params.TotalTracks, params.TotalDuration, params.UniqueArtists, params.Vibe, params.TopArtistsIds, params.TopTracksIds, params.TopAlbumsIds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user statistics", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, userStatistic)
}

func (s *Server) getUserStatisticByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	Period := r.PathValue("period")
	userStatistic, err := s.db.GetUserStatisticByPeriod(r.Context(), UserID, db.StatisticPeriod(strings.ToUpper(Period)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user statistics", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, userStatistic)
}
