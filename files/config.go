package files

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"path/filepath"
)

// Configuration is the data required to start the tool
type Configuration struct {
	Port       string
	Host       string
	DebugLevel string
	BotName    string
	UserToken  string
}

// FileExists checks if a file exists at a given path, it returns the condition
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// GenerateToken returns a random token of 50 characters
func GenerateToken() string {
	b := make([]byte, 50)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

// SetupConfig initializes the config file if it does not exists and returns the config itself
func SetupConfig(filename string) *Configuration {
	config := Configuration{}
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Dir("./"))

	if !FileExists(filename + ".toml") {

		log.Error("Config file does not exist")

		config.Host = "localhost"
		config.BotName = "Olivia"
		config.Port = "8080"
		config.DebugLevel = "error"
		config.UserToken = GenerateToken()

		viper.Set("host", config.Host)
		viper.Set("botname", config.BotName)
		viper.Set("port", config.Port)
		viper.Set("debuglevel", config.DebugLevel)
		viper.Set("usertoken", config.UserToken)

		viper.AddConfigPath(".")

		err := viper.SafeWriteConfig()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("Fatal error config file: %s \n", err)
		}

		config.Host = viper.GetString("host")
		if len(config.Host) == 0 {
			config.Host = "localhost"
		}

		config.DebugLevel = viper.GetString("debuglevel")
		if len(config.DebugLevel) == 0 {
			config.DebugLevel = "error"
		}

		config.BotName = viper.GetString("botname")
		if len(config.BotName) == 0 {
			config.BotName = "Olivia"
		}

		config.Port = viper.GetString("port")
		if len(config.Port) == 0 {
			config.Port = "8080"
		}

		config.UserToken = viper.GetString("usertoken")
		if len(config.UserToken) == 0 {
			config.UserToken = GenerateToken()
			viper.WriteConfig()
		}
	}

	return &config
}
