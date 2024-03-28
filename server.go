package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"web-server-golang/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type WebhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
			Username string `json:"username"`
		} `json:"chat"`
	} `json:"message"`
}


// This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	insert, err := db.AddInbox(body.Message.Chat.Username,  body.Message.Text)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insert)


	// Check if the message contains the word "marco"
	// if not, return without doing anything
	if !strings.Contains(strings.ToLower(body.Message.Text), "marco") {
		reqBody := &sendMessageReqBody{
			ChatID: body.Message.Chat.ID,
			Text: "Saye tak tau.. :)",
		}

		menu, err := db.ShowMenu()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(menu)
	
		// Create the JSON body from the struct
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatal(err)
		}
		
		db.AddOutbox(body.Message.Chat.Username, reqBody.Text)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post("https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			log.Fatal(res.StatusCode)
		}

		return
	}

	// If the text contains marco, call the `sayPolo` function, which
	// is defined below
	err = sayPolo(body.Message.Chat.ID, body.Message.Chat.Username)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}


type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}


func sayPolo(chatID int64, uname string) error {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   "Polo!!",
	}

	insert, err := db.AddOutbox(uname, reqBody.Text)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insert)

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}
	return nil
}


func main()  {
	var token string = "6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	db.DatabaseConnection()

	bot.Debug = true

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/webhook", http.HandlerFunc(Handler))
	
	http.ListenAndServe(":8443", nil)
	log.Println("Server is listening on localhost")
}