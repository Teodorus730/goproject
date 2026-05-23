package config

import (
    "fmt"
    "os"
	"strconv"
)

type Config struct {
    DB DBConfig
    Server ServerConfig
    RabbitMQ   RabbitMQConfig
}

type KitcherConfig struct {
	Workers  int
	RabbitMQ RabbitMQConfig
}

type RabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Vhost    string
}

type DBConfig struct {
    Host string
    Port string
    User string
    Password string
    Name string
    SSLMode string
}

type ServerConfig struct {
    Port string
    JwtSecret []byte
}

func New() (*Config, error) {
	rabbitCfg, err := LoadRabbitConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load rabbit config: %w", err)
    }

    return &Config{
        DB: DBConfig{
            Host: os.Getenv("DB_HOST"),
            Port: os.Getenv("DB_PORT"),
            User: os.Getenv("POSTGRES_USER"),
            Password: os.Getenv("POSTGRES_PASSWORD"),
            Name: os.Getenv("POSTGRES_DB"),
            SSLMode: os.Getenv("DB_SSLMODE"),
        },
        Server: ServerConfig{
            Port: os.Getenv("SERVER_PORT"),
            JwtSecret: []byte(os.Getenv("JWT_SECRET")),
        },
		RabbitMQ: rabbitCfg,
    }, nil
}

func LoadWorkersConfig() (KitcherConfig, error) {
	var cfg KitcherConfig
	var err error

	cfg.Workers, err = LoadWorkersCount()
	if err != nil {
		return KitcherConfig{}, err
	}

	cfg.RabbitMQ, err = LoadRabbitConfig()
	if err != nil {
		return KitcherConfig{}, err
	}

	return cfg, nil
}

func LoadWorkersCount() (int, error) {
	count, err := getEnvInt("WORKERS_COUNT")
	if err != nil {
		return 0, err
	}

	return count, nil
}

func LoadRabbitConfig() (RabbitMQConfig, error) {
	cfg := RabbitMQConfig{}
	var err error

	cfg.Host, err = getEnv("RABBIT_HOST")
	if err != nil {
		return RabbitMQConfig{}, err
	}

	cfg.Port, err = getEnvInt("RABBIT_PORT")
	if err != nil {
		return RabbitMQConfig{}, err
	}

	cfg.User, err = getEnv("RABBITMQ_DEFAULT_USER")
	if err != nil {
		return RabbitMQConfig{}, err
	}

	cfg.Password, err = getEnv("RABBITMQ_DEFAULT_PASS")
	if err != nil {
		return RabbitMQConfig{}, err
	}

	cfg.Vhost, err = getEnv("RABBIT_VHOST")
	if err != nil {
		return RabbitMQConfig{}, err
	}

	return cfg, nil
}

func (c *DBConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}


func getEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%s: %s", "Env not set", key)
	}
	return v, nil
}

func getEnvInt(key string) (int, error) {
	v, err := getEnv(key)
	if err != nil {
		return 0, err
	}
	intVal, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("atoi failed: %w", err)
	}
	return intVal, nil
}