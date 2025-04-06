package spotify

type Config struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://api.spotify.com/v1",
	}
}
