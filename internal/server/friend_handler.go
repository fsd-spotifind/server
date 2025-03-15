package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) addFriendHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	type parameters struct {
		FriendID string `json:"friend_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	friend, err := s.db.AddFriend(r.Context(), UserID, params.FriendID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add friend", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, friend)
}

func (s *Server) getFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	friendRequests, err := s.db.GetFriendRequests(r.Context(), UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get friend requests", err)
		return
	}
	respondWithJSON(w, http.StatusOK, friendRequests)
}

func (s *Server) acceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	type parameters struct {
		RequestID string `json:"request_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	friend, err := s.db.AcceptFriendRequest(r.Context(), UserID, params.RequestID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't accept friend", err)
		return
	}
	respondWithJSON(w, http.StatusOK, friend)
}

func (s *Server) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	UserID := r.PathValue("userId")
	friends, err := s.db.GetFriends(r.Context(), UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get friends", err)
		return
	}
	respondWithJSON(w, http.StatusOK, friends)
}
