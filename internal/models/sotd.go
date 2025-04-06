package models

type SotdEntry struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	TrackID   string `json:"trackId"`
	Note      string `json:"note"`
	Mood      string `json:"mood"`
	SetAt     string `json:"setAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Track     Track  `json:"track"`
}
