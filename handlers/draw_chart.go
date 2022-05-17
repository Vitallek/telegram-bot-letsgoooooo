package handlers

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	tele "gopkg.in/telebot.v3"
	"log"
	"tg-weather-bot-go/misc"
)

func drawChart(c tele.Context) error {
	results := misc.GetData()
	values := make([]chart.Value, 0)
	for _, r := range results {
		var c misc.City
		if err := r.Unmarshal(&c); err != nil {
			log.Fatalln("Error unmarshaling result:", err)
		}
		values = append(
			values,
			chart.Value{Label: c.City, Value: float64(c.Queries)},
		)
	}

	response := chart.BarChart{
		Title:      "Top 10 cities by queries",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Width:  512,
		Height: 512,
		XAxis:  chart.StyleShow(),
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			Range: &chart.ContinuousRange{
				Min: 0,
			},
		},
		Bars: values,
	}

	buffer := &bytes.Buffer{}
	err := response.Render(chart.PNG, buffer)
	if err != nil {
		log.Panic(err)
	}
	a := &tele.Photo{File: tele.FromReader(buffer)}
	return c.Send(a)
}
