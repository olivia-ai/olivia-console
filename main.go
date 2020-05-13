package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/olivia-ai/olivia-console/files"
	log "github.com/sirupsen/logrus"
)

const (
	logFileName    = "logfile.log"
	configFileName = "config.json"
)

var locale = "en"

// RequestMessage is the structure that uses entry connections to chat with the websocket
type RequestMessage struct {
	Type        int                    `json:"type"` // 0 for handshakes and 1 for messages
	Content     string                 `json:"content"`
	Token       string                 `json:"user_token"`
	Information map[string]interface{} `json:"information"`
	Locale      string                 `json:"locale"`
}

// ResponseMessage is the structure used to reply to the user through the websocket
type ResponseMessage struct {
	Content     string                 `json:"content"`
	Tag         string                 `json:"tag"`
	Information map[string]interface{} `json:"information"`
}

func main() {
	// Setup the logs and the config file
	files.SetupLog(logFileName)
	config := files.SetupConfig(configFileName)
	files.SetupLogLevel(*config)

	// Initialize the url
	url := url.URL{Scheme: "ws", Host: config.Host + ":" + config.Port, Path: "/websocket"}
	log.Info("Connecting to", url.String())

	// Start the connection with the websocket
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println("Unable to connect the API.")
		log.Fatal("Dial:", err)
	}
	defer c.Close()

	information := map[string]interface{}{}
	request := RequestMessage{
		Type:        0,
		Content:     "",
		Information: information,
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	if err = c.WriteMessage(websocket.TextMessage, bytes); err != nil {
		log.Error(err)
	}

	MsgType := 1

	fmt.Println(color.Magenta.Render("Enter message to " + config.BotName + " or type:"))
	fmt.Printf("- %s to quit\n", color.Green.Render("/quit"))
	fmt.Printf("- %s to change the language\n", color.Green.Render("/lang <en|fr|es...>"))
	fmt.Println()

	messagescanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">")
		messagescanner.Scan()
		text := messagescanner.Text()
		if strings.ToLower(text) == "/quit" {
			c.SetWriteDeadline(time.Now().Add(1 * time.Second))
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(1 * time.Second)
			c.Close()
			fmt.Println("Exiting from chat.")
			log.Debug("Exiting from chat")
			break
		}

		if strings.HasPrefix(text, "/lang") {
			locale = strings.Split(text, " ")[1]
			fmt.Printf("Language changed to %s.\n", color.FgMagenta.Render(locale))
			continue
		}

		secondMessage := RequestMessage{
			Type:        MsgType,
			Content:     text,
			Token:       config.UserToken,
			Information: information,
			Locale:      locale,
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
		information = response.Information
	}
}
