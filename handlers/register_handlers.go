package handlers

import tele "gopkg.in/telebot.v3"

func RegisterHandlers(b *tele.Bot) {
	b.Handle("/start", startCommand)
	b.Handle("/jh", helpCommand)
	b.Handle("/jw", weatherCommand)
	b.Handle(tele.OnLocation, location)
	b.Handle("/jt", tikTakToeCommand)
	b.Handle("/jstat", drawChart)
	b.Handle(tele.OnCallback, callbacksSwitch)
}
