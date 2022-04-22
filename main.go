package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"

	"tg-weather-bot-go/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
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
