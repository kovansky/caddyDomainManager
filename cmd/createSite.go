/*
Copyright © 2021 F4 Developer (Stanisław Kowański) <skowanski@f4dev.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/kovansky/caddyDomainManager/cmd/structs"
	"github.com/kovansky/caddyDomainManager/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	port int
)

// createSiteCmd represents the createSite command
var createSiteCmd = &cobra.Command{
	Use:   "createSite <domain name> [website type]",
	Short: "Create a new website",
	Long:  `Create a new website, including its home directory, database and user in given server (mysql, mongo) and Caddy config.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read required environment variables
		envConfig := utils.EnvironmentConfig{}

		if ok, missing := envConfig.ReadEnvironments(); !ok {
			// One of the required environment variables is missing
			println("You are missing a required environment variable ", missing)
			return
		}

		siteConfig := structs.SiteConfig{
			Type:       utils.ProgramTypePhp,
			DomainName: strings.ToLower(args[0]),
		}

		if len(args) > 1 {
			siteConfig.Type = utils.GetProgramType(args[1])
		}

		if siteConfig.Type == utils.ProgramTypeApp {
			siteConfig.Port = port
		}

		if ok, err := siteConfig.CreateConfig(envConfig); !ok {
			if os.IsExist(err) {
				println(fmt.Sprintf("Config file for domain %s already exists", siteConfig.DomainName))
				os.Exit(1)
			} else if os.IsNotExist(err) {
				println(fmt.Sprintf("Template file for type %s do not exist", strings.ToLower(string(siteConfig.Type))))
				os.Exit(1)
			} else {
				panic(err)
			}
		}

		if ok, err := siteConfig.EnableSite(envConfig); !ok {
			if os.IsNotExist(err) {
				println(fmt.Sprintf("Caddyfile for domain %s do not exist", strings.ToLower(siteConfig.DomainName)))
				os.Exit(1)
			} else {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createSiteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createSiteCmd.PersistentFlags().String("foo", "", "A help for foo")
	createSiteCmd.Flags().IntVarP(&port, "port", "p", 8080, "A port of application behind the proxy")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createSiteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
