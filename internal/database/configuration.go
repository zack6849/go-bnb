package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

// DefaultSslMode for postgres, can be disable, allow, prefer, require, verify-ca, and verify-full
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

func (c *Configuration) Open() (*sql.DB, error) {
	return sql.Open(c.Driver, c.GetConnectionString())
}

func (c *Configuration) GetDriver(db *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(db, &postgres.Config{})
}

func (c *Configuration) NewMigrator(path string, db *sql.DB) (*migrate.Migrate, error) {
	driver, err := c.GetDriver(db)
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(path, c.Driver, driver)
}
