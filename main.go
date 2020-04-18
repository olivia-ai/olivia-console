package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/olivia-ai/olivia-console/files"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	logFileName    = "logfile.log"
	configFileName = "config.json"
)

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
	// Setup the logs and the config file
	files.SetupLog(logFileName)
	config := files.SetupConfig(configFileName)
	files.SetupLogLevel(*config)

	// Initialize the url
	url := url.URL{Scheme: "ws", Host: config.Host + ":" + config.Port, Path: "/websocket"}
	log.Info("Connecting to %s..", url.String())

	// Start the connection with the websocket
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println("Unable to connect the API.")
		log.Fatal("Dial:", err)
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

