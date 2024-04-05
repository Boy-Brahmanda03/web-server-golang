package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"web-server-golang/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	bot.Debug = true

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

	//handler untuk menerima respon webhook dari API Telegram
	// http.HandleFunc("/webhook", http.HandlerFunc)
	
	//menjalankan webserver dengan goroutine (coroutine dalam go) sehingga kode dapat berjalan secara paralel
	go http.ListenAndServe(":8443", nil)
	log.Println("Bot is running and listening for updates...")


	//function yang disediakan library untuk menerima masukan update pada webhook yang dibuat
	updates := bot.ListenForWebhook("/webhook")

	//loop untuk menerima update dari req webhook
	for update := range updates {
		// Create a new MessageConfig. We don't have text yet,
        // so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message == nil {
			continue
		}

		db.AddInbox(update.Message.MessageID, int(update.Message.Chat.ID), update.Message.Chat.UserName, update.Message.Text)

		//mengambil proses state dari database
		currentState := db.GetStateMessage(update.Message.Chat.ID)
		switch currentState {
			//state awal ketika user baru masuk bot
			case 0: 
				// proses state awal
				switch update.Message.Command() {
					case "start":
						msg.Text = "Selamat datang di DailyBoyBot, ketik /help untuk melihat menu yang tersedia!"
						bot.Send(msg)
					case "help":
						msg.Text = "I understand /menu"
						bot.Send(msg)
					case "menu":
						msg.Text = menu()
						db.UpdateState(update.Message.Chat.ID, 1)
						bot.Send(msg)
					default:
						msg.Text = "Perintah tidak tersedia, silahkan gunakan /help untuk melihat menu yang tersedia!"
						bot.Send(msg)
				}
			// state ketika user telah memilih menu
			case 1:
				// proses memilih menu
				switch update.Message.Text {
					case "1":
						msg.Text = "Masukan NIM : "
						db.UpdateState(update.Message.Chat.ID, 2)
						//mengubah state pilihan menu
						db.UpdateStateMenu(update.Message.Chat.ID, 1)
						bot.Send(msg)
					case "2":
						msg.Text = "Masukan Nama Dosen : "
						db.UpdateState(update.Message.Chat.ID, 2)
						db.UpdateStateMenu(update.Message.Chat.ID, 2)
						bot.Send(msg)
					case "3":
						msg.Text = "Masukan Nama Matkul : "
						db.UpdateState(update.Message.Chat.ID, 2)
						db.UpdateStateMenu(update.Message.Chat.ID, 3)
						bot.Send(msg)
					default:
						msg.Text = "Pilihan Menu Tidak Tersedia!"
						bot.Send(msg)
				}
			//state memproses permintaan menu user
			case 2: 
				stateMenu := db.GetStateMenu(update.Message.Chat.ID)
				switch stateMenu {
				case 1:
					nim := update.Message.Text
					res, err := db.CariMahasiswa(nim)
					if err != nil {
						fmt.Println(err)
					}

					var mahasiswa []string
				
					for _, item := range res {
						mhsStrings :=  fmt.Sprintf("Mahasiswa Ditemukan!\nNIM : %s\nNama : %s\n", item.NIM, item.Nama)
						mahasiswa = append(mahasiswa, mhsStrings)
					}

					result := strings.Join(mahasiswa, "\n")
					msg.Text = result
					db.UpdateState(update.Message.Chat.ID, 0)
					bot.Send(msg)	
				case 2:
					namaDsn := update.Message.Text
					res, err := db.CariDosen(namaDsn)
					if err != nil {
						fmt.Println(err)
					}

					var dosen []string
					dosen = append(dosen, "Dosen Ditemukan!")
					for _, item := range res {
						dsnStrings :=  fmt.Sprintf("\nNIP : %s\nNIDN : %s\nNama : %s\nEmail : %s", item.NIP, item.NIDN, item.Nama, item.Email)
						dosen = append(dosen, dsnStrings)
					}
					
					result := strings.Join(dosen, "\n")
					fmt.Println("Result :", result)
					if result == "" {
						msg.Text = "Data Dosen Tidak Ditemukan"
						bot.Send(msg)
						db.UpdateState(update.Message.Chat.ID, 0)
						continue
					}
					msg.Text = result
					db.UpdateState(update.Message.Chat.ID, 0)
					bot.Send(msg)
				case 3: 
					msg.Text = "Ketemu Matkul"
					db.UpdateState(update.Message.Chat.ID, 0)
					bot.Send(msg)
				}
		}
		
		_, err := db.AddOutbox(update.Message.MessageID, int(update.Message.Chat.ID), update.Message.Chat.UserName, msg.Text)

		if err != nil {
			log.Fatalf("Error In Add Message To Outbox %s", err.Error())
		}
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