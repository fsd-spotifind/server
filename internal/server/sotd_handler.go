package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
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
	Date := r.URL.Query().Get("date")

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

	respondWithJSON(w, http.StatusOK, sotd)
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
	respondWithJSON(w, http.StatusOK, sotds)
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
