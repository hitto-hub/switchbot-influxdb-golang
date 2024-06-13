package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
    SwitchBotAPIToken string
    SwitchBotSecret   string
    InfluxDBConfig    InfluxDBConfig
}

type InfluxDBConfig struct {
    Token  string
    Org    string
    Bucket string
    URL    string
}

func LoadConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    return Config{
        SwitchBotAPIToken: os.Getenv("SWITCHBOT_API_TOKEN"),
        SwitchBotSecret:   os.Getenv("SWITCHBOT_SECRET"),
        InfluxDBConfig: InfluxDBConfig{
            Token:  os.Getenv("INFLUXDB_TOKEN"),
            Org:    os.Getenv("INFLUXDB_ORG"),
            Bucket: os.Getenv("INFLUXDB_BUCKET"),
            URL:    os.Getenv("INFLUXDB_URL"),
        },
    }
}
