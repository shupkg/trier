package jwt

type Config struct {
	SecretKey string `json:"secret_key" toml:"secret_key" yaml:"secret_key"`
	Method    string `json:"method" toml:"method" yaml:"method"`
}
