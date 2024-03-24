package main

import (
	"log"
	"net/http"
)


func main()  {

	//Url to get bot info
	//https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/getMe

	//Url To get updates
	//https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/getUpdates
	
	//Url To Send Message
	//https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/sendMessage?chat_id=1013469903&text=Hello%20User

	//URL to Set WebHook
	//https://api.telegram.org/bot6803966935:AAGIHvrBHfPVeyNFdwW9f7xouxOtFDGPPEE/setWebhook?url=https://51d2-182-253-51-29.ngrok-free.ap

	http.Handle("/", http.FileServer(http.Dir("./static")))
    log.Fatal(http.ListenAndServe(":8443", nil))
}
