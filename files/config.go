package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	log "github.com/sirupsen/logrus"
)

// Configuration is the data required to start the tool
type Configuration struct {
	Port       string `json:"port"`
	Host       string `json:"host"`
	SSL        bool   `json:"ssl"`
	DebugLevel string `json:"debug_level"`
	BotName    string `json:"bot_name"`
	UserToken  string `json:"user_token"`
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
func SetupConfig(fileName string) *Configuration {
	config := Configuration{
		Port:       "8080",
		SSL:        false,
		Host:       "localhost",
		DebugLevel: "error",
		BotName:    "Olivia",
		UserToken:  GenerateToken(),
	}

	if FileExists(fileName) {
		// Read and parse the json file
		file, _ := ioutil.ReadFile(fileName)
		err := json.Unmarshal(file, &config)
		if err != nil {
			log.Fatal(err)
		}

		return &config
	}

	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)

	return &config
}
