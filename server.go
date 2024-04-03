package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"web-server-golang/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//pilihan menu, namun belum digunakan
var menuOption = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cari Mahasiswa", ""),
    ),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cari Dosen", "2"),
    ),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cari Mata Kuliah", "3"),
    ),
)


func main() {
	//bot token
	var token string = "6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE"

	//buat object bot dengan token 
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//untuk koneksi database
	db.DatabaseConnection()

	//untuk debug respon bot
	// bot.Debug = true

	//melihat bot yang terhubung
	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	//menyambungkan webhook untuk request bot
	wh, err := tgbotapi.NewWebhook("https://busy-buck-humorous.ngrok-free.app/webhook")
	if err != nil {
		log.Fatalf("Error creating webhook: %v", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("Error setting webhook: %v", err)
	}
	
	//menjalankan webserver dengan goroutine (coroutine dalam go) sehingga kode dapat berjalan secara paralel
	go http.ListenAndServe(":8443", nil)
	log.Println("Bot is running and listening for updates...")

	//function yang disediakan library untuk menerima masukan update pada webhook yang dibuat
	updates := bot.ListenForWebhook("/webhook")
	for update := range updates {
		// Create a new MessageConfig. We don't have text yet,
        // so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)


		if update.Message == nil {
			continue
		}

		db.AddInbox(update.Message.MessageID, int(update.Message.Chat.ID), update.Message.Chat.UserName, update.Message.Text)

		if !update.Message.IsCommand() {
            msg.Text = "Pesan tidak dapat dimengerti"
			bot.Send(msg)
			addMessageToOutbox(update, msg.Text)
			continue
        }
    	
        // Extract the command from the Message.
        switch update.Message.Command() {
		case "start":
			msg.Text = "Selamat datang di DailyBoyBot, ketik /help untuk melihat menu yang tersedia!"
        case "help":
            msg.Text = "I understand /sayhi and /status."
        case "sayhi":
            msg.Text = "Hi :)"
        case "status":
            msg.Text = "I'm ok."
		case "menu":
			msg.Text = menu()
			bot.Send(msg)
			switch update.Message.CommandArguments() {
			case "cari_mhs":
				msg.Text = "Masukan NIM: "
				bot.Send(msg)
			}
        default:
            msg.Text = "I don't know that command"
			bot.Send(msg)
        }
		
		addMessageToOutbox(update, msg.Text)
	}
}

//fungsi untuk melihat menu yang dimilki dari database
func menu() string{
	menu, _ := db.ShowMenu()
	var menuStrings []string

	for _, item := range menu {
		itemStrings := fmt.Sprintf("Pilihan Menu : \n %d. %s \n> %s", item.No, item.Label, item.Deskripsi)
		menuStrings = append(menuStrings, itemStrings)
	}

	resultMenu := strings.Join(menuStrings, "\n")

	return resultMenu
}

//fungsi untuk memasukan pesan outbox ke database
func addMessageToOutbox(update tgbotapi.Update, msg string){
	_, err := db.AddOutbox(update.Message.MessageID, int(update.Message.Chat.ID), update.Message.Chat.UserName, msg)

	if err != nil {
		log.Fatalf("Error In Add Message To Outbox %s", err.Error())
	}
}