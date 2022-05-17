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
	"github.com/kovansky/caddyDomainManager/cmd/databases"
	"github.com/kovansky/caddyDomainManager/cmd/structs"
	"github.com/kovansky/caddyDomainManager/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	port            int
	forceBaseDomain bool

	dbTypeString string
	dbType       utils.DatabaseType

	dbAdminUser     string
	dbAdminPassword string
	dbHost          string
	dbAuthDatabase  string

	dbUserName     string
	dbUserPassword string
	dbUserHost     string
	dbDatabaseName string
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

		if forceBaseDomain {
			siteConfig.ForceBase = true
		}

		if ok, err := siteConfig.CreateFileStructure(envConfig); !ok {
			if os.IsNotExist(err) {
				println(fmt.Sprintf("Warning: template directory for %s type does not exist; omitting file structure copy.", strings.ToLower(string(siteConfig.Type))))
			} else if err.Error() == "domain directory not empty" {
				println(fmt.Sprintf("Warning: directory structure for %s already exists and is not empty; omitting file structure copy.", siteConfig.DomainName))
			} else {
				panic(err)
			}
		}

		println(fmt.Sprintf("[%s] Created file structure in %s using %s template", siteConfig.DomainName, siteConfig.FilesRoot(), strings.ToLower(string(siteConfig.Type))))

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

		println(fmt.Sprintf("[%s] Created Caddyfile config in %s using %s template", siteConfig.DomainName, siteConfig.Caddyfile(), strings.ToLower(string(siteConfig.Type))))

		if ok, err := siteConfig.EnableSite(envConfig); !ok {
			if os.IsNotExist(err) {
				println(fmt.Sprintf("Caddyfile for domain %s do not exist", strings.ToLower(siteConfig.DomainName)))
				os.Exit(1)
			} else {
				panic(err)
			}
		}

		println(fmt.Sprintf("[%s] Created symlink for Caddyfile in sites-enabled directory", siteConfig.DomainName))

		// Reload caddy
		siteConfig.ReloadCaddy()

		// Database configuration
		dbType = utils.GetDatabaseType(dbTypeString)
		var keysPrefix string

		switch dbType {
		case utils.DatabaseMongo:
			keysPrefix = "mongo."
			break
		case utils.DatabaseMysql:
			keysPrefix = "mysql."
			break
		case utils.DatabasePostgres:
			keysPrefix = "postgres."
			break
		}

		if dbType == utils.DatabaseNone {
			if len(dbTypeString) > 0 {
				println(fmt.Sprintf("%s is not correct type of database. Please, use 'mysql' or 'mongo'", dbTypeString))
			}
		} else {
			if len(dbAdminUser) == 0 {
				conf := viper.GetString(keysPrefix + "username")
				if len(conf) > 0 {
					dbAdminUser = conf
				} else {
					println("You are missing a database admin username (--db-admin, -U).")
					return
				}
			}

			if len(dbAdminPassword) == 0 {
				conf := viper.GetString(keysPrefix + "password")
				if len(conf) > 0 {
					dbAdminPassword = conf
				} else {
					println(fmt.Sprintf("Please, provide database password for user %s", dbAdminUser))

					bytePassword, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						panic(err)
						return
					}

					dbAdminPassword = string(bytePassword)
				}
			}

			if len(dbHost) == 0 || dbHost == "127.0.0.1" {
				conf := viper.GetString(keysPrefix + "host")
				if len(conf) > 0 && dbHost != "127.0.0.1" {
					dbHost = conf
				} else {
					var port int

					switch dbType {
					case utils.DatabaseMongo:
						port = 27017
						break
					case utils.DatabaseMysql:
						port = 3306
						break
					}

					dbHost = fmt.Sprintf("127.0.0.1:%d", port)
				}
			}

			if len(dbUserName) == 0 {
				/*
				   Build user name from domain. I.e. when domain is example.com - username is example.
				   When domain is test.example.com - username is test_example
				*/
				domainParts := siteConfig.DomainStructure(true)
				dbUserName = strings.Split(strings.Join(domainParts, "_"), ".")[0] // Get rid of TLD
			}

			if len(dbDatabaseName) == 0 {
				dbDatabaseName = dbUserName
			}

			if len(dbUserPassword) == 0 {
				dbUserPassword = utils.RandomPassword(16)
			}

			if len(dbUserHost) == 0 {
				switch dbType {
				case utils.DatabaseMysql:
					dbUserHost = "localhost"
					break
				case utils.DatabaseMongo:
					dbUserHost = "127.0.0.1"
					break
				}
			}

			splitted := strings.Split(dbHost, ":")

			if len(splitted) != 2 {
				println(fmt.Sprintf("The host (%s) is in incorrect format - it should be host:port", dbHost))
				return
			}

			dbHost = splitted[0]
			port, err := strconv.Atoi(splitted[1])

			if err != nil {
				panic(err)
			}

			// Try to create database
			var source databases.DatabaseSource
			switch dbType {
			case utils.DatabaseMongo:
				source = &databases.MongoSource{
					User:     dbAdminUser,
					Password: dbAdminPassword,
					Host:     dbHost,
					Port:     port,
					AuthDb:   dbAuthDatabase,
				}
				break
			case utils.DatabaseMysql:
				source = &databases.MongoSource{
					User:     dbAdminUser,
					Password: dbAdminPassword,
					Host:     dbHost,
					Port:     port,
				}
				break
			case utils.DatabasePostgres:
				source = &databases.PostgresSource{
					User:     dbAdminUser,
					Password: dbAdminPassword,
					Host:     dbHost,
					Port:     port,
				}
				break
			}

			if ok := source.Connect(); !ok {
				println("There was an error while connecting to the database server")
				return
			}
			defer source.Close()

			if ok := source.CreateDatabase(dbDatabaseName); !ok {
				println("There was an error while creating the database")
				return
			}

			println(fmt.Sprintf("[%s] Created database %s in %s server %s:%d", siteConfig.DomainName, dbDatabaseName, strings.ToLower(string(dbType)), dbHost, port))

			if ok := source.CreateUser(dbUserName, dbUserHost, dbUserPassword); !ok {
				println("There was an error while creating the database user")
				return
			}

			siteConfig.WriteDatabaseInfo(dbHost, port, dbDatabaseName, dbUserName, dbUserPassword, dbUserHost)

			println(fmt.Sprintf("[%s] Created user %s (with connection limited to %s) and granted privileges on %s in %s server %s:%d. All required information were stored in database_info.txt file in website's root directory", siteConfig.DomainName, dbUserName, dbUserHost, dbDatabaseName, strings.ToLower(string(dbType)), dbHost, port))
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
	createSiteCmd.Flags().BoolVarP(&forceBaseDomain, "basedomain", "b", false, "Force to treat the domain as high-level, even if contains subdomains")

	createSiteCmd.Flags().StringVarP(&dbTypeString, "db-type", "t", "", "Type of database to use (MySQL, Postgres or Mongo). If this flag is present an user in corresponding database will be created. Requires providing all other database-related flags.")

	createSiteCmd.Flags().StringVarP(&dbAdminUser, "db-admin", "U", "", "Database administrator username")
	createSiteCmd.Flags().StringVarP(&dbAdminPassword, "db-admin-password", "P", "", "Database administrator password")
	createSiteCmd.Flags().StringVarP(&dbHost, "db-host", "H", "127.0.0.1", "Database hostname (with port)")
	createSiteCmd.Flags().StringVarP(&dbAuthDatabase, "db-auth-db", "s", "", "Authentication database (only for mongo)")

	createSiteCmd.Flags().StringVarP(&dbUserName, "username", "u", "", "Name of the database user to create. Optional, default extracted from the domain name.")
	createSiteCmd.Flags().StringVarP(&dbUserPassword, "password", "i", "", "Password of the database user to create. Optional, randomly generated by default.")
	createSiteCmd.Flags().StringVarP(&dbUserHost, "host", "o", "", "Host to which the database user to create should be limited while connecting. Optional, localhost by default.")
	createSiteCmd.Flags().StringVarP(&dbDatabaseName, "database", "D", "", "Name of the database to create. Optional, default equal to the username.")

	viper.BindPFlag("mongo.authDatabase", createSiteCmd.Flag("db-auth-db"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createSiteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
