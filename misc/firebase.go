package misc

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
)

type City struct {
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	Queries int    `json:"queries,omitempty"`
}

func SaveData(cityID int, city string, region string, country string, queries int) error {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://telegramweatherdb-default-rtdb.europe-west1.firebasedatabase.app",
	}
	opt := option.WithCredentialsFile("firebase-admin-key.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	ref := client.NewRef("/cities/")
	results, err := ref.OrderByChild("queries").GetOrdered(ctx)
	if err != nil {
		log.Fatalln("Error querying database:", err)
	}
	for _, r := range results {
		var c City
		if err := r.Unmarshal(&c); err != nil {
			log.Fatalln("Error unmarshaling result:", err)
		}

		if r.Key() == fmt.Sprintf("%d", cityID) {
			cityQueriesToChange := c.Queries
			return ref.Child(fmt.Sprintf("%d", cityID)).Set(ctx, City{
				City:    city,
				Region:  region,
				Country: country,
				Queries: cityQueriesToChange + 1,
			})
		}
	}

	return ref.Child(fmt.Sprintf("%d", cityID)).Set(ctx, City{
		City:    city,
		Region:  region,
		Country: country,
		Queries: queries,
	})
}

func GetData() []db.QueryNode {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://telegramweatherdb-default-rtdb.europe-west1.firebasedatabase.app",
	}
	opt := option.WithCredentialsFile("firebase-admin-key.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	ref := client.NewRef("/cities/")

	results, err := ref.OrderByChild("queries").LimitToLast(10).GetOrdered(ctx)
	if err != nil {
		log.Fatalln("Error querying database:", err)
	}
	return results
}
