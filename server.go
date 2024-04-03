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

// BAGIAN REQUEST BOT
// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type WebhookReqBody struct {
	Message struct {
		MessageID int64 `json:"message_id"`
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

	//menyimpan pesan masuk ke dalam database
	_, err := db.AddInbox(int(body.Message.MessageID),int(body.Message.Chat.ID) ,body.Message.Chat.Username,  body.Message.Text)
	if err != nil {
		log.Fatalf("Error creating webhook: %v", err)
		return
	}


	// Check if the message contains the word "menu"
	if strings.Contains(strings.ToLower(body.Message.Text), "menu") {
		err := menuResponse(body)
		if err != nil {
			log.Fatalf("Error creating webhook: %v", err)
			return
		}
		return
	}


	if (strings.ToLower(body.Message.Text) == "1") {
		reqBody := &sendMessageReqBody{
			ChatID: body.Message.Chat.ID,
			MessageID: body.Message.MessageID,
			Text:   "Menu Cari Mahasiswa\nMasukan NIM : ",
		}

		db.AddOutbox(int(body.Message.MessageID),int(body.Message.Chat.ID), body.Message.Chat.Username, reqBody.Text)

		// Create the JSON body from the struct
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			return 
		}
	
		// Send a post request with your token
		res, err := http.Post("https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			return
		}
		if res.StatusCode != http.StatusOK {
			return 
		} 
		return
	}


	if (strings.ToLower(body.Message.Text) != "") {
		if (strings.ToLower(body.Message.Text) != "1") && len(body.Message.Text) == 10 {
			nim := body.Message.Text
			res, err := db.CariMahasiswa(nim)
			if err != nil {
				log.Println("Error querying database:", err)
				return
			}
			if res != nil {
				reqBody := &sendMessageReqBody{
					ChatID: body.Message.Chat.ID,
					MessageID: body.Message.MessageID,
					Text: "",
				}

				var mahasiswa []string
				
				for _, item := range res {
					mhsStrings :=  fmt.Sprintf("Mahasiswa dengan NIM %s:\nNama: %s\n", item.NIM, item.Nama)
					mahasiswa = append(mahasiswa, mhsStrings)
				}
			
				result := strings.Join(mahasiswa, "\n")
				reqBody.Text = result

				// Create the JSON body from the struct
				reqBytes, err := json.Marshal(reqBody)
				if err != nil {
					log.Fatal(err)
				}

				db.AddOutbox(int(body.Message.MessageID), int(body.Message.Chat.ID), body.Message.Chat.Username, reqBody.Text)
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
			} else {
				reqBody := &sendMessageReqBody{
					ChatID: body.Message.Chat.ID,
					MessageID: body.Message.MessageID,
					Text: "Data NIM tidak tersedia",
				}
				db.AddOutbox(int(body.Message.MessageID), int(body.Message.Chat.ID), body.Message.Chat.Username, reqBody.Text)
				if err != nil {
					log.Fatal(err)
				}

				// Create the JSON body from the struct
				reqBytes, err := json.Marshal(reqBody)
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
		}
	}


	// If the text not contains on the list, call the `defaultResponse` function, which
	err = defaultResponse(body)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}
}


//		BAGIAN RESPON BOT
type sendMessageReqBody struct {
	MessageID int64 `json:"reply_to_message_id"`
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func defaultResponse(body *WebhookReqBody) error {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: body.Message.Chat.ID,
		MessageID: body.Message.MessageID,
		Text:   "Menu Tidak Tersedia.\nSilahkan ketik menu untuk melihat menu yang tersedia!",
	}

	db.AddOutbox(int(body.Message.MessageID),int(body.Message.Chat.ID), body.Message.Chat.Username, reqBody.Text)

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

func menuResponse(body *WebhookReqBody) error {
	reqBody := &sendMessageReqBody{
		ChatID: body.Message.Chat.ID,
		MessageID: body.Message.MessageID,
		Text: "",
	}

	menu, _ := db.ShowMenu()
	var menuStrings []string
	
	for _, item := range menu {
		itemStrings := fmt.Sprintf("Menu Item \n %d. %s \n> %s", item.No, item.Label, item.Deskripsi)
		menuStrings = append(menuStrings, itemStrings)
	}
	
	result := strings.Join(menuStrings, "\n")
	reqBody.Text = result

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	
	db.AddOutbox(int(body.Message.MessageID), int(body.Message.Chat.ID), body.Message.Chat.Username, reqBody.Text)
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

	return nil
}


func main() {
	var token string = "6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	db.DatabaseConnection()

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	wh, err := tgbotapi.NewWebhook("https://busy-buck-humorous.ngrok-free.app/webhook")
	if err != nil {
		log.Fatalf("Error creating webhook: %v", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("Error setting webhook: %v", err)
	}

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.HandleFunc("/webhook", http.HandlerFunc(Handler))
	
	//menjalankan webserver dengan goroutine (coroutine dalam go) sehingga kode dapat berjalan secara paralel
	go http.ListenAndServe(":8443", nil)
	log.Println("Bot is running and listening for updates...")
	updates := bot.ListenForWebhook("/webhook")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }
     
        // Create a new MessageConfig. We don't have text yet,
        // so we leave it empty.
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

        // Extract the command from the Message.
        switch update.Message.Command() {
        case "help":
            msg.Text = "I understand /sayhi and /status."
        case "sayhi":
            msg.Text = "Hi :)"
        case "status":
            msg.Text = "I'm ok."
        default:
            msg.Text = "I don't know that command"
        }

        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
	}
}