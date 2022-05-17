package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type PostgresSource struct {
	User     string
	Password string
	Host     string
	Port     int

	db           *sql.DB
	databaseName string
}

func (source *PostgresSource) Connect() bool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d", source.User, source.Password, source.Host, source.Port)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false
	}

	source.db = db

	return true
}

func (source PostgresSource) CreateUser(name string, _ string, password string) bool {
	_, err := source.db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", name, password))
	if err != nil {
		println("1", err.Error())
		return false
	}

	_, err = source.db.Exec(fmt.Sprintf("GRANT ALL ON %s TO %s", source.databaseName, name))
	if err != nil {
		println("3", err.Error())
		return false
	}

	return true
}

func (source *PostgresSource) CreateDatabase(name string) bool {
	_, err := source.db.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		return false
	}

	source.databaseName = name

	return true
}

func (source PostgresSource) Close() {
	_ = source.db.Close()
}
