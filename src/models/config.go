package models

type Config struct {
	SpotifyData      SpotifyData
	DbProvider       string
	ConnectionString string
}

type SpotifyData struct {
	AuthAPI  string
	TrackAPI string
	ClientID string
	Secret   string
}

func NewConfig() Config {
	return Config{
		ConnectionString: "server=PAVILLALOBOS;user id=;trusted_connection=true;database=spotify;app name=spotify",
		DbProvider:       "mssql",
		SpotifyData: SpotifyData{
			AuthAPI:  "https://accounts.spotify.com/api/token",
			TrackAPI: "https://api.spotify.com/v1/tracks/%s",
			ClientID: "f12f4ec67fa34ddb97840aa122a9e331",
			Secret:   "cd9b981f1cbf4f4393a20ddf6ebafb74",
		},
	}
}
