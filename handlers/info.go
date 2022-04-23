package handlers

import tele "gopkg.in/telebot.v3"

func startCommand(ctx tele.Context) error {
	return ctx.Send("Hello from Yura, senior Go developer\n/jh - best command")
}

func helpCommand(ctx tele.Context) error {
	return ctx.Send("/jh - список команд\n/jw [город] - погода\n/jt - крестики-нолики\n/jstat - статистика запросов городов")
}
