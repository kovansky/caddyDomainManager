package structs

import (
	"fmt"
	"github.com/kovansky/caddyDomainManager/cmd/utils"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type CaddyConfig struct {
	Type       utils.ProgramType
	DomainName string
	Port       int
}

func (c CaddyConfig) CreateConfig(envConfig utils.EnvironmentConfig) (bool, error) {
	caddyfileNameFormat := "%s.Caddyfile"

	// Set locations
	sitesAllPath := path.Join(envConfig.CaddySites, "sites-all")

	// Get template path
	templateName := fmt.Sprintf("template_%s", strings.ToLower(string(c.Type)))
	templatePath := path.Join(sitesAllPath, templateName)

	// Create Caddyfile path
	caddyfileDestinationName := fmt.Sprintf(caddyfileNameFormat, c.DomainName)
	destinationPath := path.Join(sitesAllPath, caddyfileDestinationName)

	if !fileExists(templatePath) {
		return false, fs.ErrNotExist
	}

	// Check, if Caddyfile for this domain do not already exist
	if fileExists(destinationPath) {
		return false, fs.ErrExist
	}

	// Read template
	template, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return false, nil
	}

	// Replace vars in template
	templateSpecific := strings.ReplaceAll(string(template), "$SITE_ADDRESS", c.DomainName)
	templateSpecific = strings.ReplaceAll(templateSpecific, "$PORT", strconv.Itoa(c.Port))

	// Write Caddyfile
	err = ioutil.WriteFile(destinationPath, []byte(templateSpecific), 0775)
	if err != nil {
		return false, err
	}

	return true, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
