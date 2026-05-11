package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8080"`

	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`

	JWTSecret          string        `envconfig:"JWT_SECRET"           required:"true"`
	AccessTokenExpiry  time.Duration `envconfig:"ACCESS_TOKEN_EXPIRY"  default:"15m"`
	RefreshTokenExpiry time.Duration `envconfig:"REFRESH_TOKEN_EXPIRY" default:"168h"`
	BCryptCost         int           `envconfig:"BCRYPT_COST"          default:"12"`

	Env string `envconfig:"ENV" default:"development"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
