package server

import (
	"encoding/json"
	"fsd-backend/internal/auth"
	"fsd-backend/internal/logger"
	"fsd-backend/prisma/db"
	"net/http"
	"strings"
	"time"
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
		logger.ErrorLog.Printf("Couldn't hash password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	user, err := s.db.CreateUser(r.Context(), params.Email, params.Username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "Unique constraint failed") {
			respondWithError(w, http.StatusBadRequest, "Email or username already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User         *db.UserModel
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.ErrorLog.Printf("Couldn't decode parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	user, err := s.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		s.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	err = s.db.CreateRefreshToken(
		r.Context(),
		refreshToken,
		user.ID,
		time.Now().Add(time.Hour*24*60), // 60 days
	)
	if err != nil {
		logger.ErrorLog.Printf("Couldn't store refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})

}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find bearer token", err)
		return
	}

	user, err := s.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		s.jwtSecret,
		time.Hour,
	)
	if err != nil {
		logger.ErrorLog.Printf("Couldn't create access token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (s *Server) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find bearer token", err)
		return
	}

	err = s.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		logger.ErrorLog.Printf("Couldn't revoke refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find bearer token", err)
		return
	}

	err = s.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't logout", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
