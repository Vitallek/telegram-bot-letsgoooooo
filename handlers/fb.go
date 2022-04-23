package handlers

import (
	firebase "firebase.google.com/go"
	//"github.com/buger/jsonparser"
	//"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
)

func SaveData(cityID int,city string,region string,country string,queries int) error {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://telegramweatherdb-default-rtdb.europe-west1.firebaseio.com",
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("firebase-admin-key.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("/cities/")

	// User is a json-serializable type.
	type City struct {
		city    string `json:"city,omitempty"`
		region  string `json:"region,omitempty"`
		country string `json:"country,omitempty"`
		queries int    `json:"queries,omitempty"`
	}

	err = ref.Set(ctx, map[string]*City{
		string(rune(cityID)): {
			city:    city,
			region:  region,
			country: country,
			queries: queries,
		},
	})
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
	return nil
}

