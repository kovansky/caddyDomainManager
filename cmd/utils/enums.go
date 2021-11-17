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
	case "APP":
	case "APPLICATION":
	case "PROXY":
		return ProgramTypeApp
	case "PHP":
	case "WP":
	case "WORDPRESS":
		return ProgramTypePhp
	case "HTML":
	default:
		return ProgramTypeHtml
	}

	return ProgramTypeHtml
}
