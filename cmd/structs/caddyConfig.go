package structs

import "github.com/kovansky/caddyDomainManager/cmd/utils"

type CaddyConfig struct {
	Type       utils.ProgramType
	DomainName string
	Port       int
}

func (c CaddyConfig) createConfig() {

}
