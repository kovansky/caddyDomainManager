package databases

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlSource struct {
	User     string
	Password string
	Host     string
	Port     int

	db           *sql.DB
	databaseName string
}

func (source *MysqlSource) Connect() bool {
	dsn := mysql.Config{
		User:              source.User,
		Passwd:            source.Password,
		Addr:              fmt.Sprintf("%s:%d", source.Host, source.Port),
		InterpolateParams: true,
	}

	db, err := sql.Open("mysql", dsn.FormatDSN())

	if err != nil {
		return false
	}

	err = db.Ping()
	if err != nil {
		return false
	}

	source.db = db

	return true
}

func (source MysqlSource) CreateUser(name string, userHost string, password string) bool {
	_, err := source.db.Exec(fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY '%s'", name, userHost, password))
	if err != nil {
		println("1", err.Error())
		return false
	}

	_, err = source.db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%s'", source.databaseName, name, userHost))
	if err != nil {
		println("3", err.Error())
		return false
	}

	flushPrivilegesStmt, err := source.db.Prepare("FLUSH PRIVILEGES")
	if err != nil {
		println("5", err.Error())
		return false
	}

	_, err = flushPrivilegesStmt.Exec()
	if err != nil {
		println("6", err.Error())
		return false
	}

	return true
}

func (source *MysqlSource) CreateDatabase(name string) bool {
	_, err := source.db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	if err != nil {
		return false
	}

	source.databaseName = name

	return true
}

func (source MysqlSource) Close() {
	_ = source.db.Close()
}
