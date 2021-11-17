package structs

import (
	"errors"
	"fmt"
	"github.com/kovansky/caddyDomainManager/cmd/utils"
	copyDirs "github.com/otiai10/copy"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type SiteConfig struct {
	Type       utils.ProgramType
	DomainName string
	Port       int
	ForceBase  bool
	caddyfile  string
	filesRoot  string
}

// Functions regarding Caddy

func (cfg *SiteConfig) CreateConfig(envConfig utils.EnvironmentConfig) (bool, error) {
	caddyfileNameFormat := "%s.Caddyfile"

	// Set locations
	sitesAllPath := path.Join(envConfig.CaddySites, "sites-all")

	// Get template path
	templateName := fmt.Sprintf("template_%s", strings.ToLower(string(cfg.Type)))
	templatePath := path.Join(sitesAllPath, templateName)

	// Create Caddyfile path
	caddyfileDestinationName := fmt.Sprintf(caddyfileNameFormat, cfg.DomainName)
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
	templateSpecific := strings.ReplaceAll(string(template), "$SITE_ADDRESS", cfg.DomainName)
	templateSpecific = strings.ReplaceAll(templateSpecific, "$FILES_ROOT", cfg.filesRoot)
	templateSpecific = strings.ReplaceAll(templateSpecific, "$PORT", strconv.Itoa(cfg.Port))

	// Write Caddyfile
	err = ioutil.WriteFile(destinationPath, []byte(templateSpecific), 0775)
	if err != nil {
		return false, err
	}

	cfg.caddyfile = destinationPath

	return true, nil
}

func (cfg SiteConfig) EnableSite(envConfig utils.EnvironmentConfig) (bool, error) {
	// Set locations
	sitesEnabledPath := path.Join(envConfig.CaddySites, "sites-enabled")
	fileName := filepath.Base(cfg.caddyfile)

	// Check, if Caddyfile for this domain exists
	if fileExists(cfg.caddyfile) {
		return false, fs.ErrNotExist
	}

	// Create symlink in sites-enabled
	err := os.Symlink(cfg.caddyfile, path.Join(sitesEnabledPath, fileName))

	if err != nil {
		return false, err
	}

	return true, nil
}

// Functions regarding file structure

func (cfg *SiteConfig) CreateFileStructure(envConfig utils.EnvironmentConfig) (bool, error) {
	// Set locations
	templatesBasePath := path.Join(envConfig.ServerFiles, "templates")
	templatePath := path.Join(templatesBasePath, strings.ToLower(string(cfg.Type)))

	if !directoryExists(templatePath) {
		return false, fs.ErrNotExist
	}

	// Create target path
	domainStructure := ReverseSlice(cfg.DomainStructure())
	domainRootPath := ""
	currentDomain := ""

	for index, domain := range domainStructure {
		currentDomain = fmt.Sprintf("%s.%s", domain, currentDomain)
		addPath := currentDomain

		if index > 0 {
			addPath = path.Join("domains", addPath)
		}

		domainRootPath = path.Join(domainRootPath, addPath)
	}

	domainRootPath = path.Join(envConfig.ServerFiles, domainRootPath)

	if directoryExists(domainRootPath) {
		dir, err := os.Open(domainRootPath)
		if err != nil {
			return false, err
		}
		defer func(dir *os.File) {
			_ = dir.Close()
		}(dir)

		_, err = dir.Readdirnames(1)
		if err == io.EOF {
			return false, errors.New("domain directory not empty")
		}
	}

	err := os.MkdirAll(domainRootPath, 0775)
	if err != nil {
		return false, err
	}

	cfg.filesRoot = domainRootPath

	err = copyDirs.Copy(templatePath, domainRootPath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cfg SiteConfig) DomainStructure() []string {
	if cfg.ForceBase {
		return []string{cfg.DomainName}
	}

	splitted := strings.Split(cfg.DomainName, ".")

	if len(splitted) == 1 {
		return splitted
	}
	if len(splitted) == 0 {
		return []string{}
	}

	splitted[len(splitted)-2] = splitted[len(splitted)-2] + "." + splitted[len(splitted)-1]
	splitted = splitted[:len(splitted)-1]

	return splitted
}

func ReverseSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}
