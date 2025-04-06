package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"fsd-backend/internal/database"
	"fsd-backend/internal/spotify"
)

type Server struct {
	port    int
	db      database.Service
	spotify spotify.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	spotifyConfig := spotify.DefaultConfig()
	spotifyConfig.ClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyConfig.ClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	NewServer := &Server{
		port:    port,
		db:      database.New(),
		spotify: spotify.NewService(spotifyConfig),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
