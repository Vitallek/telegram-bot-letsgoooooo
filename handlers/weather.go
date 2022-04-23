package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	tele "gopkg.in/telebot.v3"

	"tg-weather-bot-go/misc"
)

func weatherCommand(ctx tele.Context) error {
	message := ctx.Message().Text
	messageArray := strings.Fields(message)
	if len(messageArray) == 1 {
		return ctx.Send("Формат ввода:\n/jw [город]")
	}
	cityName := strings.Join(messageArray[1:], " ")

	langCode := "ru"
	if strings.ContainsAny("abcdefghigklmnopqrstuvwxyz", strings.ToLower(string(cityName[0]))) {
		langCode = "en"
	}

	url := fmt.Sprintf("https://wft-geo-db.p.rapidapi.com/v1/geo/cities?limit=5&namePrefix=%s&sort=-population&languageCode=%s&types=CITY", cityName, langCode)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Host", "wft-geo-db.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("GEODB_TOKEN"))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return ctx.Send("Город не найден")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic()
	}

	kb := &tele.ReplyMarkup{}
	btnRowArray := []tele.Row{}
	jsonparser.ArrayEach(responseData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetInt(value, "id")
		city, _ := jsonparser.GetString(value, "city")
		region, _ := jsonparser.GetString(value, "region")
		country, _ := jsonparser.GetString(value, "country")
		fullCityName := misc.GetFullCityName(city, region, country, ", ")

		btnWeather := kb.Data(fullCityName, fmt.Sprintf("w_%d", id))
		btnRowArray = append(btnRowArray, kb.Row(btnWeather))
	}, "data")
	kb.Inline(btnRowArray...)

	return ctx.Send("Найденные города: \n", kb)
}

func weatherCallback(ctx tele.Context) error {
	callbackArray := strings.Split(ctx.Callback().Data, "_")
	cityId := callbackArray[1]

	url := fmt.Sprintf("https://wft-geo-db.p.rapidapi.com/v1/geo/cities/%s?languageCode=ru", cityId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Host", "wft-geo-db.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("GEODB_TOKEN"))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return ctx.Send("Город не найден")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic()
	}

	cityID, _ := jsonparser.GetInt(responseData, "data", "id")
	lat, _ := jsonparser.GetFloat(responseData, "data", "latitude")
	lon, _ := jsonparser.GetFloat(responseData, "data", "longitude")
	city, _ := jsonparser.GetString(responseData, "data", "city")
	region, _ := jsonparser.GetString(responseData, "data", "region")
	country, _ := jsonparser.GetString(responseData, "data", "country")
	fullCityName := misc.GetFullCityName(city, region, country, ", ")
	latStr := fmt.Sprintf("%f", lat)
	lonStr := fmt.Sprintf("%f", lon)
	log.Print(cityID)

	err = SaveData(int(cityID),city,region,country,1)
	if err != nil {
		return err
	}
	
	err = ctx.Send(misc.GetWeather(latStr, lonStr, fullCityName))
	if err != nil {
		return err
	}
	return ctx.Respond()
}

func location(ctx tele.Context) error {
	lat := ctx.Message().Location.Lat
	lon := ctx.Message().Location.Lng
	return ctx.Send(misc.GetWeather(fmt.Sprintf("%f", lat), fmt.Sprintf("%f", lon), ""))
}
