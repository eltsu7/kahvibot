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

func kirjaus(update tgbotapi.Update, kuvaus string) {

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Muuttujat updatesta
	var aika int = update.Message.Date
	var userID int = update.Message.From.ID
	var updateID int = update.UpdateID
	if kuvaus == "" {
		kuvaus = update.Message.CommandArguments()
	}
	var nimi string = update.Message.From.UserName

	// Lisää uusi kupillinen
	sql := "insert into juonnit (update_id, user_id, aika, kuvaus) values ($1, $2, to_timestamp($3), $4)"

	rows, err := conn.Query(context.Background(), sql, updateID, userID, aika, kuvaus)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	// Päivitä nimi nimitauluun
	sql = "INSERT INTO nimet VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username"

	rows, err = conn.Query(context.Background(), sql, userID, nimi)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()
}

func santsi(update tgbotapi.Update) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Hae viimeisimmän kirjatun kupin kuvaus
	var kuvaus string
	sql := "SELECT kuvaus FROM juonnit WHERE user_id=$1 ORDER BY id DESC LIMIT 1;"
	err = conn.QueryRow(context.Background(), sql, update.Message.From.ID).Scan(&kuvaus)
	if err == pgx.ErrNoRows {
		return "Ei edellistä kuppia."
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	// Lisää uusi kupillinen vanhalla kuvauksella
	kirjaus(update, kuvaus)
	log.Println("Santsi:", update.Message.From.UserName, kuvaus)
	return ""
}

func kupit(id int) string {

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Laske kupilliset ja palauta lukumäärä
	var kuppeja int
	sql := "select kupit from kuppilaskuri where user_id=$1"
	err = conn.QueryRow(context.Background(), sql, id).Scan(&kuppeja)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var txt string
	if kuppeja == 0 {
		txt = fmt.Sprint("Et taida olla kovin sivistynyt ihminen.")
	} else if kuppeja == 1 {
		txt = fmt.Sprint("Olet juonut yhden kupin kahvia. Tsemppaatko kiitos.")
	} else {
		txt = fmt.Sprint("Olet juonut ", kuppeja, " kuppia kahvia.")
	}

	return txt
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

			// Muuttujat updatesta
			userID := update.Message.From.ID
			userName := update.Message.From.UserName
			chatID := update.Message.Chat.ID

			switch update.Message.Command() {
			case "help":
				txt := "Hei minä olen kahvikanalabotti.\n\nKomennot:\n" +
					"/kahvi [valinnainen kuvaus]\nKirjaa uuden kupillisen.\n" +
					"/santsi\nKirjaa kupillisen, mutta kopioi kuvauksen sun edellisestä kupista.\n" +
					"/kupit\nKertoo montako kuppia oot juonu"
				msgHelp := tgbotapi.NewMessage(chatID, txt)
				bot.Send(msgHelp)
				log.Println("Help:", userName)

			case "kahvi":
				var kuvaus string = update.Message.CommandArguments()
				if len(kuvaus) > 255 {
					liianPitkaMsg := tgbotapi.NewMessage(chatID, "Liian pitkä kuvaus..")
					bot.Send(liianPitkaMsg)
				} else {
					kirjaus(update, "")
					log.Println("Kirjaus:", userName, kuvaus)
				}

			case "santsi":
				var err string = santsi(update)
				if err != "" {
					log.Println("Santsi kusee, ei edellistä kuppia", userName)
					errMsg := tgbotapi.NewMessage(chatID, err)
					bot.Send(errMsg)
				}

			case "kupit":
				txt := kupit(userID)

				msgStats := tgbotapi.NewMessage(chatID, txt)
				bot.Send(msgStats)
				log.Println("Kupit:", userName)
			}
		}

	}
}
