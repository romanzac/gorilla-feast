package config

type Config struct {
	Web struct {
		Listen     string
		Port       string
		DisableTLS string
		Key        string
		Cert       string
		JWTPrivKey string
		JWTPubKey  string
	}
	Database struct {
		PostgresURI string
	}
}

var (
	Cfg Config
)
