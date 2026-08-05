package configuration

import (
	"GoBNB/internal/database"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetString(key, def string) string {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return raw
}

func GetInt(key string, def int) int {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		fmt.Printf("config: failed to parse value %s through parse function for key %s", raw, key)
		return def
	}
	return parsed
}

func GetDatabaseConfiguration() *database.Configuration {

	return &database.Configuration{
		Driver:   GetString(database.DriverKey, database.DefaultDriver),
		Username: GetString(database.UsernameKey, database.DefaultUsername),
		Password: GetString(database.PasswordKey, database.DefaultPassword),
		Hostname: GetString(database.HostnameKey, database.DefaultHostname),
		Port:     GetInt(database.PortKey, database.DefaultPort),
		Schema:   GetString(database.SchemaKey, database.DefaultSchema),
		SSLMode:  GetString(database.SslKey, database.DefaultSslMode),
	}
}
