package handlers

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
)

func callbacksSwitch(ctx tele.Context) error {
	callback := ctx.Callback().Data
	switch strings.Split(callback, "_")[0] {
	case "\fw":
		return weatherCallback(ctx)
	case "\ftreg":
		return tikTakToeRegCallback(ctx)
	case "\ftplay":
		return tikTakToePlaceMarkCallback(ctx)
	case "\ftdead":
		return tikTakToeDeadCallback(ctx)
	}
	log.Panic(callback)
	return nil
}
