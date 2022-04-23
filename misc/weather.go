package misc

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
)

func GetWeather(lat string, lon string, cityName string) string {
	response, err := http.Get(fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&exclude=minutely,hourly,alerts&appid=%s&units=metric&lang=ru",
		lat, lon, os.Getenv("OPEN_WEATHER_TOKEN")))
	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	println(string(responseData))

	weatherDesc, err := jsonparser.GetString(responseData, "current", "weather", "[0]", "description")
	temperature, err := jsonparser.GetFloat(responseData, "current", "temp")
	feelsLike, err := jsonparser.GetFloat(responseData, "current", "feels_like")
	humidity, err := jsonparser.GetFloat(responseData, "current", "humidity")
	wind, err := jsonparser.GetFloat(responseData, "current", "wind_speed")

	sendTemplate := "Погода"
	if cityName != "" {
		sendTemplate += " в городе " + cityName
	}
	sendTemplate += ":\n\nОблачность: " + weatherDesc + "\n" +
		"Температура: " + fmt.Sprintf("%d", int(math.Round(temperature))) + " °C\n" +
		"Ощущается как " + fmt.Sprintf("%d", int(math.Round(feelsLike))) + " °C\n" +
		"Влажность: " + fmt.Sprintf("%d", int(math.Round(humidity))) + "%\n" +
		"Ветер:" + fmt.Sprintf("%d", int(math.Round(wind))) + " м/c"

	return sendTemplate
}
