package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Read the config file named 'pf_dev' and return the configuration
func ReadConfig() Configuration {

	viper.SetConfigName("pf_dev")
	viper.AddConfigPath("./")
	viper.SetConfigType("yml")

	var Config Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return Config
}

// Get the HTTP server configuration
func GetConfig() (string, string) {

	httpConfig := ReadConfig()

	httpHost := httpConfig.Server.Host
	httpPort := fmt.Sprintf(":%d", httpConfig.Server.Port)

	return httpHost, httpPort

}

// Get the Pprof server configuration
func GetPprofConfig() (bool, string, string) {

	pprofConfig := ReadConfig()

	pprofEnable := pprofConfig.PprofServer.Enable
	pprorHost := pprofConfig.PprofServer.Host
	pprofPort := fmt.Sprintf(":%d", pprofConfig.PprofServer.Port)

	return pprofEnable, pprorHost, pprofPort

}

// Get the client configuration
func GetClientConfig() (string, string) {

	clientConfig := ReadConfig()

	clientHost := clientConfig.Client.Host
	clientPort := fmt.Sprintf(":%d", clientConfig.Client.Port)

	return clientHost, clientPort

}
