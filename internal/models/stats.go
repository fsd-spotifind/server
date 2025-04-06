package models

import (
	db "fsd-backend/prisma/db"
)

type ComputedStats struct {
	TotalTracks   int
	TotalDuration int
	UniqueArtists int
	Vibe          string
	TopArtistIDs  []string
	TopTrackIDs   []string
	TopAlbumIDs   []string
}

type UserStatistic struct {
	ID            string             `json:"id"`
	UserID        string             `json:"userId"`
	Period        db.StatisticPeriod `json:"period"`
	TotalTracks   int                `json:"totalTracks"`
	TotalDuration int                `json:"totalDuration"`
	UniqueArtists int                `json:"uniqueArtists"`
	Vibe          string             `json:"vibe"`
	TopTracks     []Track            `json:"topTracks"`
	TopArtists    []Artist           `json:"topArtists"`
	TopAlbums     []Album            `json:"topAlbums"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
}
