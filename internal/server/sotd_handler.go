package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) createSotdHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID  string `json:"user_id"`
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

	sotd, err := s.db.CreateSotd(r.Context(), params.UserID, params.TrackID, params.Note, params.Mood)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create Sotd", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, sotd)
}

func (s *Server) getSotdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	unixTime, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Unix timestamp format", http.StatusBadRequest)
		return
	}

	parsedDate := time.Unix(unixTime, 0).UTC()

	sotd, err := s.db.GetSotd(r.Context(), userId, parsedDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get Sotd", err)
		return
	}

	respondWithJSON(w, http.StatusOK, sotd)
}

func (s *Server) updateSotdHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		SotdID  string `json:"sotd_id"`
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

	sotd, err := s.db.UpdateSotd(r.Context(), params.SotdID, params.TrackID, params.Note, params.Mood)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update Sotd", err)
		return
	}

	respondWithJSON(w, http.StatusOK, sotd)
}
