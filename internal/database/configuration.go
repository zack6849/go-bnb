package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DriverKey = "DB_DRIVER"
const DefaultDriver = "postgres"
const UsernameKey = "DB_USERNAME"
const DefaultUsername = "gobnb"
const PasswordKey = "DB_PASSWORD"
const DefaultPassword = "secret"
const HostnameKey = "DB_HOSTNAME"
const DefaultHostname = "localhost"
const PortKey = "DB_PORT"
const DefaultPort = 3306
const SchemaKey = "DB_SCHEMA"
const DefaultSchema = "gobnb"
const SslKey = "DB_SSL"

// DefaultSslMode for postgres, can be one of: disable, allow, prefer, require, verify-ca, and verify-full
// https://www.postgresql.org/docs/current/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS
const DefaultSslMode = "prefer"

type Configuration struct {
	Driver   string
	Username string
	Password string
	Hostname string
	Port     int
	Schema   string
	SSLMode  string
}

func (c *Configuration) GetConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		c.Driver,
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.Schema,
		c.SSLMode,
	)
}

func (c *Configuration) Open() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.GetConnectionString()), &gorm.Config{})
}
