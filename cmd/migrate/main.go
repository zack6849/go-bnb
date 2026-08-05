package main

import (
	"GoBNB/internal/bootstrap"
	"GoBNB/internal/configuration"
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//setup the application
	bootstrap.Initialize()
	dbConfig := configuration.GetDatabaseConfiguration()
	conn, err := dbConfig.Open()
	if err != nil {
		fmt.Printf("failed to connect to DB %s: %s\n", dbConfig.GetConnectionString(), err.Error())
		return
	}
	//ensure we close the db connection once it's opened
	defer cleanupConnection(conn)
	migrator, err := dbConfig.NewMigrator("file://db/migrations", conn)
	if err != nil {
		fmt.Printf("failed to initialize migrator: %s\n", err.Error())
		return
	}
	err = migrator.Up()
	if err != nil {
		fmt.Printf("db migrations failed: %s\n", err.Error())
	}
}

func cleanupConnection(conn *sql.DB) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("failed to close database connection: %s\n", err.Error())
	}
}
