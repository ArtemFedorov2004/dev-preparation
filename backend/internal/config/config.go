package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Keycloak KeycloakConfig
	Redis    RedisConfig
	Swagger  SwaggerConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type KeycloakConfig struct {
	URL             string
	Realm           string
	RefreshInterval time.Duration
}

type RedisConfig struct {
	Addr        string
	Password    string
	DB          int
	TopicTTL    time.Duration
	QuestionTTL time.Duration
}

type SwaggerConfig struct {
	OAuth2AuthorizationUrl string
	OAuth2TokenUrl         string
}

func (c KeycloakConfig) RealmURL() string {
	return fmt.Sprintf("%s/realms/%s",
		c.URL, c.Realm,
	)
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	keycloak := KeycloakConfig{
		URL:             getEnvStr("KEYCLOAK_URL", "http://localhost:8443"),
		Realm:           getEnvStr("KEYCLOAK_REALM", "devprep"),
		RefreshInterval: getEnvDuration("KEYCLOAK_JWKS_REFRESH", 15*time.Minute),
	}

	oidcBase := keycloak.RealmURL() + "/protocol/openid-connect"

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8081),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnvStr("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			Name:            getEnvStr("DB_NAME", "devprep"),
			User:            getEnvStr("DB_USER", "devprep"),
			Password:        getEnvStr("DB_PASSWORD", "devprep"),
			SSLMode:         getEnvStr("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Keycloak: keycloak,
		Redis: RedisConfig{
			Addr:        getEnvStr("REDIS_ADDR", "localhost:6379"),
			Password:    getEnvStr("REDIS_PASSWORD", ""),
			DB:          getEnvInt("REDIS_DB", 0),
			TopicTTL:    getEnvDuration("REDIS_TOPIC_TTL", 10*time.Minute),
			QuestionTTL: getEnvDuration("REDIS_QUESTION_TTL", 10*time.Minute),
		},
		Swagger: SwaggerConfig{
			OAuth2AuthorizationUrl: getEnvStr("SWAGGER_OAUTH2_AUTHORIZATION_URL", oidcBase+"/auth"),
			OAuth2TokenUrl:         getEnvStr("SWAGGER_OAUTH2_TOKEN_URL", oidcBase+"/token"),
		},
	}

	return cfg, nil
}

func getEnvStr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}
