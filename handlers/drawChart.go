package handlers

import (
	"bytes"
	"github.com/wcharczuk/go-chart" //exposes "chart"
	tele "gopkg.in/telebot.v3"
	"log"
)

func draw(c tele.Context) error{
	// Get values from the DB and convert
	results := GetData()
	values := make([]chart.Value, 0)
	for _, r := range results {
		var c City
		if err := r.Unmarshal(&c); err != nil {
			log.Fatalln("Error unmarshaling result:", err)
		}
		//your code here
		//fmt.Printf("%s has %d queries\n", c.City, c.Queries)
		values = append(
			values,
			chart.Value{Label: c.City, Value: float64(c.Queries)},
		)
	}

	// Chart settings
	response := chart.BarChart{
		Title:      "Top 10 cities by queries",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Width:    512,
		Height:   512,
		BarWidth: 50,
		XAxis:    chart.StyleShow(),
		YAxis: chart.YAxis{
			Style:          chart.StyleShow(),
			//ValueFormatter: chart.IntValueFormatter,
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: values[0].Value,
			},
		},
		Bars: values,
	}

	// Render and send
	buffer := &bytes.Buffer{}
	image := response.Render(chart.PNG, buffer)
	
	return c.Send(imageToSend)
}
