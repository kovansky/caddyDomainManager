package utils

import "strings"

type ProgramType string

const (
	ProgramTypeApp  ProgramType = "APPLICATION"
	ProgramTypePhp              = "PHP"
	ProgramTypeHtml             = "HTML"
)

func GetProgramType(s string) ProgramType {
	s = strings.ToUpper(s)
	switch s {
	case "APP", "APPLICATION", "PROXY":
		return ProgramTypeApp
	case "PHP", "WP", "WORDPRESS":
		return ProgramTypePhp
	default:
		return ProgramTypeHtml
	}
}

type DatabaseType string

const (
	DatabaseMysql DatabaseType = "MYSQL"
	DatabaseMongo              = "MONGO"
	DatabaseNone               = ""
)

func GetDatabaseType(s string) DatabaseType {
	s = strings.ToUpper(s)
	switch s {
	case "MONGO":
		return DatabaseMongo
	case "MYSQL":
		return DatabaseMysql
	default:
		return DatabaseNone
	}
}
