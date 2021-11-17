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
