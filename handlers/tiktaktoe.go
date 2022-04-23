package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tg-weather-bot-go/misc"

	tele "gopkg.in/telebot.v3"
)

func tikTakToeCommand(ctx tele.Context) error {
	kb := &tele.ReplyMarkup{}
	kb.Inline(kb.Row(kb.Data("Вступить в игру", fmt.Sprintf("treg_%d", ctx.Message().Sender.ID))))
	return ctx.Send("Крестики нолики", kb)
}

func tikTakToeRegCallback(ctx tele.Context) error {
	callbackData := ctx.Callback().Data
	firstPlayerId := strings.Split(callbackData, "_")[1]
	secondPlayerId := fmt.Sprintf("%d", ctx.Callback().Sender.ID)
	if firstPlayerId == secondPlayerId && os.Getenv("DEBUG") == "False" {
		return ctx.Respond()
	}

	kb := &tele.ReplyMarkup{}
	kbRows := []tele.Row{}
	for i := 0; i < 3; i++ {
		btnsArray := []tele.Btn{}
		for j := 0; j < 3; j++ {
			iStr := fmt.Sprintf("%d", i)
			jStr := fmt.Sprintf("%d", j)
			btnsArray = append(btnsArray,
				kb.Data("□", fmt.Sprintf("tplay_%s_%s_1_%s_%s_0_0_0_0_0_0_0_0_0", firstPlayerId, secondPlayerId, iStr, jStr)))
		}
		kbRows = append(kbRows, kb.Row(btnsArray...))
	}
	kb.Inline(kbRows...)

	err := ctx.Edit("Крестики нолики\nХодит первый игрок", kb)
	if err != nil {
		return err
	}
	return ctx.Respond()
}

func tikTakToePlaceMarkCallback(ctx tele.Context) error {
	callbackDataArray := strings.Split(ctx.Callback().Data, "_")
	firstPlayerId := callbackDataArray[1]
	secondPlayerId := callbackDataArray[2]
	currentPlayer := callbackDataArray[3]
	iiStr := callbackDataArray[4]
	jjStr := callbackDataArray[5]
	ii, _ := strconv.Atoi(iiStr)
	jj, _ := strconv.Atoi(jjStr)

	var field [][]string
	for i := 0; i < 3; i++ {
		field = append(field, callbackDataArray[6+3*i:9+3*i])
	}

	if currentPlayer == "1" && firstPlayerId != fmt.Sprintf("%d", ctx.Callback().Sender.ID) && os.Getenv("DEBUG") == "False" {
		return ctx.Respond()
	}
	if currentPlayer == "2" && secondPlayerId != fmt.Sprintf("%d", ctx.Callback().Sender.ID) && os.Getenv("DEBUG") == "False" {
		return ctx.Respond()
	}
	if field[ii][jj] != "0" {
		return ctx.Respond()
	}
	if currentPlayer == "1" {
		field[ii][jj] = "1"
	} else {
		field[ii][jj] = "2"
	}

	nextPlayer := "2"
	if currentPlayer == "2" {
		nextPlayer = "1"
	}

	message := "Крестики нолики"
	if nextPlayer == "1" {
		message += "\nХод первого игрока"
	} else {
		message += "\nХод второго игрока"
	}

	callback_left := fmt.Sprintf("tplay_%s_%s_%s_", firstPlayerId, secondPlayerId, nextPlayer)
	callback_right := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			callback_right += "_" + field[i][j]
		}
	}

	winner := misc.TikTakToeCheckForWin(field)
	if winner != 0 {
		message = "Крестики нолики\n"
		if winner == 1 {
			message += "Победил первый игрок"
			callback_left = "tdead_"
			callback_right = ""
		}
		if winner == 2 {
			message += "Победил второй игрок"
			callback_left = "tdead_"
			callback_right = ""
		}
		if winner == 3 {
			message += "Ничья"
			callback_left = "tdead_"
			callback_right = ""
		}
	}

	kb := &tele.ReplyMarkup{}
	kbRows := []tele.Row{}
	for i := 0; i < 3; i++ {
		btnsArray := []tele.Btn{}
		for j := 0; j < 3; j++ {
			iStr := fmt.Sprintf("%d", i)
			jStr := fmt.Sprintf("%d", j)
			currentSymbol := "□"
			if field[i][j] == "1" {
				currentSymbol = "\u274C"
			}
			if field[i][j] == "2" {
				currentSymbol = "\u2B55"
			}
			callback := fmt.Sprintf("%s%s_%s%s", callback_left, iStr, jStr, callback_right)
			btnsArray = append(btnsArray, kb.Data(currentSymbol, callback))
		}
		kbRows = append(kbRows, kb.Row(btnsArray...))
	}
	kb.Inline(kbRows...)

	err := ctx.Edit(message, kb)
	if err != nil {
		return err
	}
	return ctx.Respond()
}

func tikTakToeDeadCallback(ctx tele.Context) error {
	return ctx.Respond()
}
