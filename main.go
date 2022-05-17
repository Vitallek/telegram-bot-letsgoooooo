package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"tg-weather-bot-go/handlers"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Файл .env отсутствует")
	}

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	handlers.RegisterHandlers(b)

	b.Start()
}
