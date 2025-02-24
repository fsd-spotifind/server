package server

import (
	"encoding/json"
	"fsd-backend/internal/auth"
	"net/http"
	"strings"
)

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := s.db.CreateUser(r.Context(), params.Email, params.Username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "Unique constraint failed") {
			respondWithError(w, http.StatusBadRequest, "Email or username already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
