package models

type SpotifyUser struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Product     string    `json:"product"`
	Country     string    `json:"country"`
	Images      []Image   `json:"images"`
	Followers   Followers `json:"followers"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Followers struct {
	Total int `json:"total"`
}

type Track struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Artists    []ArtistSimplified `json:"artists"`
	Album      Album              `json:"album"`
	DurationMs int                `json:"duration_ms"`
	Popularity int                `json:"popularity"`
}

type Artist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Genres     []string `json:"genres"`
	Images     []Image  `json:"images"`
	Popularity int      `json:"popularity"`
}

type ArtistSimplified struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	TotalTracks int     `json:"total_tracks"`
	Images      []Image `json:"images"`
	ReleaseDate string  `json:"release_date"`
}

type Playlist struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      []Image        `json:"images"`
	Owner       SpotifyUser    `json:"owner"`
	Tracks      PlaylistTracks `json:"tracks"`
}

type PlaylistTracks struct {
	Total int     `json:"total"`
	Items []Track `json:"items"`
}

type RecentlyPlayed struct {
	Items   []PlayHistory `json:"items"`
	Total   int           `json:"total"`
	Next    string        `json:"next"`
	Limit   int           `json:"limit"`
	Cursors Cursors       `json:"cursors"`
}

type PlayHistory struct {
	PlayedAt string `json:"played_at"`
	Track    Track  `json:"track"`
}

type Cursors struct {
	After  string `json:"after"`
	Before string `json:"before"`
}
