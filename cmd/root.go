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
	"github.com/spf13/cobra"
	"os"
	"path"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "caddyDomainManager",
	Aliases: []string{"cdm"},
	Short:   "A CLI app to manage your Caddy Server configuration",
	Long:    `A CLI application that lets you easily create new Caddy sites from existing templates, speeding up the process of adding new domains`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.caddyDomainManager.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cdm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cdm")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		sampleViper := viper.New()

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		sampleViper.AddConfigPath(home)
		sampleViper.SetConfigType("yaml")
		sampleViper.SetConfigName(".cdm.sample")

		sampleViper.Set("mysql", map[string]string{
			"host":     "",
			"username": "",
			"password": "",
		})

		sampleViper.Set("mongo", map[string]string{
			"host":         "",
			"username":     "",
			"password":     "",
			"authDatabase": "",
		})

		err = sampleViper.SafeWriteConfig()
		if err == nil {
			println(fmt.Sprintf("Created sample config file at %s. You can fill it and rename to .cdm.yaml to make it work.", path.Join(home, ".cdm.sample.yaml")))
		}
	}
}
