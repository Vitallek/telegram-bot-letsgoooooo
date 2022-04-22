package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	tele "gopkg.in/telebot.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(ctx tele.Context) error {
		return ctx.Send("Hello from fucking GO language, motherfucker")
	})

	b.Handle("/jew", func(ctx tele.Context) error {
		msg := ctx.Message().Text
		cityNameSplitArray := strings.Fields(msg)
		cityName := cityNameSplitArray[1]
		fmt.Print(cityName + "\n")
		responseCities, err := http.Get(`https://api.geoapify.com/v1/geocode/autocomplete?text=` +
			cityName +
			`&type=city&limit=5&apiKey=` + os.Getenv("geoapifyToken"))

		if err != nil {
			fmt.Print(err.Error())
			return ctx.Send("No cities found")
		}

		responseCitiesData, err := ioutil.ReadAll(responseCities.Body)
		if err != nil {
			log.Panic()
		}

		selector := &tele.ReplyMarkup{}
		btnRawArray := []tele.Row{}
		btnArray := []tele.Btn{}
		//fmt.Print(string(responseCitiesData))
		jsonparser.ArrayEach(responseCitiesData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			//fmt.Print(string(value))
			lat, err := jsonparser.GetFloat(value, "properties", "lat")
			lon, err := jsonparser.GetFloat(value, "properties", "lon")
			cityName, err := jsonparser.GetString(value, "properties", "city")
			if cityName == "" {
				cityName = "city"
			}
			latstr := strconv.FormatFloat(float64(lat), 'f', 10, 64)
			lonstr := strconv.FormatFloat(float64(lon), 'f', 10, 64)
			btnWeather := selector.Data(cityName, `lat=`+latstr+`&lon=`+lonstr)
			btnRawArray = append(btnRawArray, selector.Row(btnWeather))
			btnArray = append(btnArray, btnWeather)

		}, "features")

		selector.Inline(
			btnRawArray...,
		)
		return ctx.Send("Я нашёл несколько городов: \n", selector)
	})
	b.Handle(tele.OnCallback, func(c tele.Context) error {
		coords := c.Callback().Data
		responseWeather, err := http.Get(`https://api.openweathermap.org/data/2.5/weather?` +
			string(coords)[2:] +
			`&limit=1&units=metric&lang=ru&appid=` + os.Getenv("openWeatherToken"))

		if err != nil {
			fmt.Print(err.Error())
		}

		responseWeatherData, err := ioutil.ReadAll(responseWeather.Body)
		if err != nil {
			log.Fatal(err)
		}

		//latTmplt, err := jsonparser.GetFloat(responseWeatherData, "coord", "lat")
		//lonTmplt, err := jsonparser.GetFloat(responseWeatherData, "coord", "lon")
		cityNameTmplt, err := jsonparser.GetString(responseWeatherData, "name")
		weatherDescTmplt, err := jsonparser.GetString(responseWeatherData, "weather", "[0]", "description")
		tempTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "temp")
		feelsTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "feels_like")
		humidityTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "humidity")
		windTmplt, err := jsonparser.GetFloat(responseWeatherData, "wind", "speed")

		sendTemplate := "\n\nПогода в городе " + cityNameTmplt + " сейчас:\n\n" +
			"Облачность: " + weatherDescTmplt +
			"\nТемпература: " + strconv.FormatFloat(tempTmplt, 'f', 0, 64) +
			"°C \nощущается как " +
			strconv.FormatFloat(feelsTmplt, 'f', 0, 64) + " °C\n" +
			"\nВлажность: " + strconv.FormatFloat(humidityTmplt, 'f', 0, 64) + "%" +
			"\nВетер:" + strconv.FormatFloat(windTmplt, 'f', 0, 64) + "м/c"

		return c.Send(string(sendTemplate))
	})
	b.Handle(tele.OnLocation, func(ctx tele.Context) error {
		lat := ctx.Message().Location.Lat
		lon := ctx.Message().Location.Lng

		//responseCity, err := http.Get(`https://api.openweathermap.org/geo/1.0/reverse?lat=` +
		//	strconv.FormatFloat(float64(lat), 'f', 10, 64) +
		//	`&lon=` +
		//	strconv.FormatFloat(float64(lon), 'f', 10, 64) +
		//	`&limit=1&appid=` + os.Getenv("openWeatherToken"))
		//
		//if err != nil {
		//	fmt.Print(err.Error())
		//	os.Exit(1)
		//}
		//
		//responseCityData, err := ioutil.ReadAll(responseCity.Body)
		//if err != nil {
		//	log.Fatal(err)
		//}

		//cityName, err := jsonparser.GetString(responseData, "[0]", "name")

		responseWeather, err := http.Get(`https://api.openweathermap.org/data/2.5/weather?lat=` +
			strconv.FormatFloat(float64(lat), 'f', 10, 64) +
			`&lon=` +
			strconv.FormatFloat(float64(lon), 'f', 10, 64) +
			`&limit=1&units=metric&lang=ru&appid=` + os.Getenv("openWeatherToken"))

		if err != nil {
			fmt.Print(err.Error())
		}

		responseWeatherData, err := ioutil.ReadAll(responseWeather.Body)
		if err != nil {
			log.Fatal(err)
		}

		latTmplt, err := jsonparser.GetFloat(responseWeatherData, "coord", "lat")
		lonTmplt, err := jsonparser.GetFloat(responseWeatherData, "coord", "lon")
		cityNameTmplt, err := jsonparser.GetString(responseWeatherData, "name")
		weatherDescTmplt, err := jsonparser.GetString(responseWeatherData, "weather", "[0]", "description")
		tempTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "temp")
		feelsTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "feels_like")
		humidityTmplt, err := jsonparser.GetFloat(responseWeatherData, "main", "humidity")
		windTmplt, err := jsonparser.GetFloat(responseWeatherData, "wind", "speed")

		sendTemplate := "Ваши координаты: \n[" + strconv.FormatFloat(latTmplt, 'f', 6, 64) + "," +
			strconv.FormatFloat(lonTmplt, 'f', 6, 64) + "] " +
			"\n\nБлижайший населённый пункт - " + string(cityNameTmplt) +
			"\n\nПогода по координатам сейчас:\n\n" +
			"Облачность: " + weatherDescTmplt +
			"\nТемпература: " + strconv.FormatFloat(tempTmplt, 'f', 2, 64) +
			"°C \nощущается как " +
			strconv.FormatFloat(feelsTmplt, 'f', 2, 64) + " °C\n" +
			"\nВлажность: " + strconv.FormatFloat(humidityTmplt, 'f', 0, 64) + "%" +
			"\nВетер:" + strconv.FormatFloat(windTmplt, 'f', 2, 64) + "м/c"

		return ctx.Send(string(sendTemplate))
	})

	b.Start()

}
