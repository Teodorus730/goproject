package config

import (
    "fmt"
    "os"
)

type Config struct {
    DB DBConfig
    Server ServerConfig
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
    }, nil
}

func (c *DBConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}
