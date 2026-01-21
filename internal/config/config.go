package config

import (
	"go.uber.org/zap"
)

type Config struct {
	ENV            string `yaml:"env" json:"env"`
	Host           string `yaml:"host" json:"host"`
	Port           string `yaml:"port" json:"port"`
	JWTSecret      string `yaml:"jwt_secret" json:"jwt_secret"`
	UserSvcAddr    string `yaml:"user_svc_addr" json:"user_svc_addr"`
	UserSvcPort    string `yaml:"user_svc_port" json:"user_svc_port"`
	PrivateKeyPath string `yaml:"private_key_path" json:"private_key_path"`
	PublicKeyPath  string `yaml:"public_key_path" json:"public_key_path"`
	JWTKeyID       string `yaml:"jwt_key_id" json:"jwt_key_id"`
	JWTIssuer      string `yaml:"jwt_issuer" json:"jwt_issuer"`
	JWTAudience    string `yaml:"jwt_audience" json:"jwt_audience"`
	Logger         *zap.Logger
}

// ConfigLoader defines the interface for loading configuration
type ConfigLoader interface {
	Load(cfg *Config) error
}

// Cfg is the global configuration instance
var Cfg *Config

// init loads the global config at startup
func init() {
	loader := &EnvConfigLoader{} // default .env path
	Cfg = loader.Load()
	logger := NewLogger(Cfg.ENV)
	Cfg.Logger = logger
}
