package utils

import "os"

type EnvironmentConfig struct {
	CaddySites      string
	ServerFiles     string
	GlobalCaddyfile string
}

func (envConfig *EnvironmentConfig) ReadEnvironments() (bool, string) {
	ok := false

	if ok, envConfig.CaddySites = getEnvNotEmpty("CADDY_SITES_DIR"); !ok {
		return false, "CADDY_SITES_DIR"
	}

	if ok, envConfig.CaddySites = getEnvNotEmpty("CADDY_GLOBAL"); !ok {
		return false, "CADDY_GLOBAL"
	}

	if ok, envConfig.ServerFiles = getEnvNotEmpty("SERVER_FILES_DIR"); !ok {
		return false, "SERVER_FILES_DIR"
	}

	return true, ""
}

func getEnvNotEmpty(key string) (bool, string) {
	if env := os.Getenv(key); env == "" {
		return false, ""
	} else {
		return true, env
	}
}
