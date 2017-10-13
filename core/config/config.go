package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load the config file and pull in ENV vars.
func Load() error {

	if isCi := os.Getenv("CI"); isCi == "true" {
		viper.SetConfigName("config")
	} else {
		viper.SetConfigName("config")
	}

	viper.AddConfigPath("/go/src/github.com/teamwork/teamwork-contacts/") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/go/src/github.com/teamwork/teamwork-contacts/")
	viper.AddConfigPath("/go/src/github.com/teamwork/teamwork-contacts/") // call multiple times to add many search paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return viper.ReadInConfig()

}
