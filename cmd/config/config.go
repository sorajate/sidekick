/*
Copyright © 2024 Mahmoud Mousa <m.mousa@hey.com>

Licensed under the GNU GPL License, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"os"

	"github.com/mightymoud/sidekick/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Subcommands for modifying sidekick config",
}

var currentContextCmd = &cobra.Command{
	Use:   "current",
	Short: "Get the current context from sidekick config",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := utils.GetSidekickConfigFromCmdContext(cmd)
		if err != nil {
			pterm.Fatal.Println(err)
		}
		fmt.Println(config.CurrentContext)
	},
}

var useContextCmd = &cobra.Command{
	Use:   "use [context-name]",
	Short: "Switch the current context in sidekick config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := utils.GetSidekickConfigFromCmdContext(cmd)
		if err != nil {
			pterm.Fatal.Println(err)
		}

		contextName := args[0]
		config.CurrentContext = contextName
		config.Save(viper.GetString("config"))
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the old sidekick config to the new config",
	Run: func(cmd *cobra.Command, args []string) {
		var serverConfig utils.SidekickServer
		content, err := os.ReadFile(viper.GetString("config"))
		if err != nil {
			pterm.Fatal.Println(err)
		}
		err = yaml.Unmarshal(content, &serverConfig)
		if err != nil {
			pterm.Fatal.Println(err)
		}

		serverConfig.Name = "default"
		defaultContext := utils.SidekickContext{
			Name:   "default",
			Server: "default",
		}
		newConfig := utils.SidekickConfig{
			Version:        "1",
			Servers:        []utils.SidekickServer{serverConfig},
			Contexts:       []utils.SidekickContext{defaultContext},
			CurrentContext: defaultContext.Name,
		}
		err = newConfig.Print()
		if err != nil {
			pterm.Fatal.Println(err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(currentContextCmd)
	ConfigCmd.AddCommand(useContextCmd)
	ConfigCmd.AddCommand(migrateCmd)
}
