package main

import (
	"TelegramBot/app"
	"fmt"
	"github.com/yanzay/tbot/v2"
)

type BotApp struct {
	client *tbot.Client
}

func init() {
	//e := godotenv.Load()
	//if e != nil {
	//	log.Println(e)
	//}
	// TODO /donote
	// Class of note
	// Name of note
	// Text of note
	// Deadline
	// в storage создается json file

	//token = "2137641983:AAFHzmAuUXsD63PDChxszpDho0TNktJyoCI"
}

func main() {
	bot := app.NewBot()
	err := bot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
//
//func (a *BotApp) andHandler ( message *tbot.Message) {
//	msg := "This is first message"
//	a.client.SendMessage(message.Chat.ID,msg)
//}