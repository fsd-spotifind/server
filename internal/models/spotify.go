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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Artists     []Artist `json:"artists"`
	Album       Album    `json:"album"`
	DurationMs  int      `json:"duration_ms"`
	Popularity  int      `json:"popularity"`
	PreviewURL  string   `json:"preview_url"`
	ExternalURL string   `json:"external_urls.spotify"`
}

type Artist struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Images      []Image `json:"images"`
	Popularity  int     `json:"popularity"`
	ExternalURL string  `json:"external_urls.spotify"`
}

type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Images      []Image `json:"images"`
	ReleaseDate string  `json:"release_date"`
	ExternalURL string  `json:"external_urls.spotify"`
}

type Playlist struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      []Image        `json:"images"`
	Owner       SpotifyUser    `json:"owner"`
	Tracks      PlaylistTracks `json:"tracks"`
	ExternalURL string         `json:"external_urls.spotify"`
}

type PlaylistTracks struct {
	Total int     `json:"total"`
	Items []Track `json:"items"`
}

type RecentlyPlayed struct {
	Items []PlayHistory `json:"items"`
	Next  string        `json:"next"`
}

type PlayHistory struct {
	PlayedAt string `json:"played_at"`
	Track    Track  `json:"track"`
}
