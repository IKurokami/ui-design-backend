package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	Env                     string
	ClientOrigin            string
	DBDriver                string
	DBDSN                   string
	JWTSecret               string
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	RefreshTokenRememberTTL time.Duration
	RotationGracePeriod     time.Duration
	SMTPHost                string
	SMTPPort                int
	SMTPUser                string
	SMTPPass                string
	MailFrom                string
	AppBaseURL              string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")
	clientOrigin := getEnv("CLIENT_ORIGIN", "http://localhost:5173")
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbDSN := getEnv("DB_DSN", "data/auth.db")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-jwt-signing-key-change-in-production-12345")

	accessTTL := parseDuration(getEnv("ACCESS_TOKEN_TTL", "15m"), 15*time.Minute)
	refreshTTL := parseDuration(getEnv("REFRESH_TOKEN_TTL", "168h"), 7*24*time.Hour)
	refreshRememberTTL := parseDuration(getEnv("REFRESH_TOKEN_REMEMBER_TTL", "720h"), 30*24*time.Hour)
	gracePeriod := parseDuration(getEnv("ROTATION_GRACE_PERIOD", "30s"), 30*time.Second)

	smtpHost := getEnv("SMTP_HOST", "localhost")
	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "1025"))
	if err != nil {
		smtpPort = 1025
	}
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASS", "")
	mailFrom := getEnv("MAIL_FROM", "Mộc <no-reply@uidesign.local>")
	appBaseURL := getEnv("APP_BASE_URL", "http://localhost:5173")

	return &Config{
		Port:                    port,
		Env:                     env,
		ClientOrigin:            clientOrigin,
		DBDriver:                dbDriver,
		DBDSN:                   dbDSN,
		JWTSecret:               jwtSecret,
		AccessTokenTTL:          accessTTL,
		RefreshTokenTTL:         refreshTTL,
		RefreshTokenRememberTTL: refreshRememberTTL,
		RotationGracePeriod:     gracePeriod,
		SMTPHost:                smtpHost,
		SMTPPort:                smtpPort,
		SMTPUser:                smtpUser,
		SMTPPass:                smtpPass,
		MailFrom:                mailFrom,
		AppBaseURL:              appBaseURL,
	}, nil
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func parseDuration(val string, defaultVal time.Duration) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
