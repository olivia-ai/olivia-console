package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	logFileName    = "logfile.log"
	configFileName = "config"
)

type Configuration struct {
	Port       string
	Host       string
	DebugLevel string
	BotName    string
	UserToken  string
}

// RequestMessage is the structure that uses entry connections to chat with the websocket
type RequestMessage struct {
	Type        int         `json:"type"` // 0 for handshakes and 1 for messages
	Content     string      `json:"content"`
	Token       string      `json:"user_token"`
	Information Information `json:"information"`
}

// ResponseMessage is the structure used to reply to the user through the websocket
type ResponseMessage struct {
	Content     string      `json:"content"`
	Tag         string      `json:"tag"`
	Information Information `json:"information"`
}

// Information is the user's information retrieved from the client
type Information struct {
	Name           string        `json:"name"`
	MovieGenres    []string      `json:"movie_genres"`
	MovieBlacklist []string      `json:"movie_blacklist"`
	Reminders      []Reminder    `json:"reminders"`
	SpotifyToken   *oauth2.Token `json:"spotify_token"`
	SpotifyID      string        `json:"spotify_id"`
	SpotifySecret  string        `json:"spotify_secret"`
}

type Reminder struct {
	Reason string `json:"reason"`
	Date   string `json:"date"`
}

func main() {

	setupLog(logFileName)
	config := setupConfig(configFileName)
	setupLogLevel(*config)

	u := url.URL{Scheme: "ws", Host: config.Host + ":" + config.Port, Path: "/websocket"}
	log.Info("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	inf := Information{
		Name: "",
	}

	request := RequestMessage{
		Type:        0,
		Content:     "",
		Information: inf,
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	if err = c.WriteMessage(websocket.TextMessage, bytes); err != nil {
		log.Error(err)
	}

	MsgType := 1

	fmt.Println(color.FgRed.Render("Enter message to " + config.BotName + " (for finish - type 'quit'):"))

	messagescanner := bufio.NewScanner(os.Stdin)

	for {
		messagescanner.Scan()
		text := messagescanner.Text()
		if strings.ToLower(text) == "quit" {
			c.SetWriteDeadline(time.Now().Add(1 * time.Second))
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(1 * time.Second)
			c.Close()
			fmt.Println("Exiting from chat.")
			log.Debug("Exiting from chat")
			break
		}

		secondMessage := RequestMessage{
			Type:        MsgType,
			Content:     text,
			Token:       config.UserToken,
			Information: inf,
		}

		newbytes, err := json.Marshal(secondMessage)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Debug("Marshaled message: " + string(newbytes))
		}

		// Send message to server
		if err = c.WriteMessage(websocket.TextMessage, newbytes); err != nil {
			log.Error(err)
		} else {
			log.Debug("Message was sended")
		}

		// Read message from server
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		} else {
			log.Debug("Get message: " + string(msg))
		}

		MsgType = msgType

		// Unmarshal the json content of the message
		var response ResponseMessage
		if err = json.Unmarshal(msg, &response); err != nil {
			log.Debug(err)
			continue
		}

		fmt.Println(color.FgYellow.Render(config.BotName + "> " + response.Content))
		inf = response.Information
	}
}
func setupLogLevel(configuration Configuration) {

	switch configuration.DebugLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}

func setupConfig(filename string) *Configuration {

	config := Configuration{}
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Dir("./"))

	if !isExists(filename + ".toml") {

		log.Error("Config file does not exist")

		config.Host = "localhost"
		config.BotName = "Olivia"
		config.Port = "8080"
		config.DebugLevel = "error"
		config.UserToken = generateUserToken(200)

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
			config.UserToken = generateUserToken(200)
			viper.WriteConfig()
		}
	}

	return &config
}

func generateUserToken(length int) string {

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_)°^¨$*£ù%=+:/;.,?-(}{[]&é@#"
	b := make([]rune, length)
	for i := range b {
		b[i] = rune(chars[rand.Intn(len(chars))])
	}
	return string(b)
}

func isExists(path string) bool {

	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func setupLog(filename string) {

	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)
	// You can change the Timestamp format. But you have to use the same date and time.
	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}
}
