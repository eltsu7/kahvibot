package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

func kirjaus(id int, kahvilaji string, aika int) {

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Lisää uusi kupillinen
	var query1 string
	if kahvilaji == "" {
		query1 = fmt.Sprint("INSERT INTO juonnit (user_id, aika) VALUES (", id, ", to_timestamp(", aika, "));")
	} else {
		query1 = fmt.Sprint("INSERT INTO juonnit (user_id, aika, kahvi) VALUES (", id, ", to_timestamp(", aika, "), '", kahvilaji, "');")
	}
	rows, err1 := conn.Query(context.Background(), query1)
	if err1 != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err1)
		os.Exit(1)
	}
	rows.Close()
}

func kupit(id int) int {

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Laske kupilliset ja palauta lukumäärä
	var kuppeja int
	query := fmt.Sprint("select count(*) from juonnit where user_id=", id)
	err = conn.QueryRow(context.Background(), query).Scan(&kuppeja)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return kuppeja
}

func main() {
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Panic(err1)
	}

	var tgToken string = os.Getenv("TG_TOKEN")

	bot, err2 := tgbotapi.NewBotAPI(tgToken)
	if err2 != nil {
		log.Panic(err2)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u) // todo: mitä jos tää heittää errorii??

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {

			userID := update.Message.From.ID
			chatID := update.Message.Chat.ID

			switch update.Message.Command() {
			case "kahvi":
				var kahvilaatu string = update.Message.CommandArguments()
				if len(kahvilaatu) > 30 {
					liianPitkaMsg := tgbotapi.NewMessage(chatID, "Liian pitkä nimi kahvillas..")
					bot.Send(liianPitkaMsg)
				} else {
					aika := update.Message.Date
					kirjaus(userID, kahvilaatu, aika)
					log.Println("Kirjaus:", userID, kahvilaatu)
				}
			case "kupit":
				kupit := kupit(userID)
				var txt string
				if kupit == 0 {
					txt = fmt.Sprint("Et taida olla kovin sivistynyt ihminen.")
				} else if kupit == 1 {
					txt = fmt.Sprint("Olet juonut yhden kupin kahvia. Tsemppaatko kiitos.")
				} else {
					txt = fmt.Sprint("Olet juonut ", kupit, " kuppia kahvia.")
				}
				msgStats := tgbotapi.NewMessage(chatID, txt)
				bot.Send(msgStats)
				log.Println("Kupit", userID)
			}
		}

	}
}
