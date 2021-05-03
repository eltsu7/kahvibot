package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func kirjaus(update tgbotapi.Update, kuvaus string) {
	var aika int = update.Message.Date
	var userID int = update.Message.From.ID
	if kuvaus == "" {
		kuvaus = update.Message.CommandArguments()
	}
	var nimi string = update.Message.From.UserName

	dbKirjaus(userID, aika, kuvaus, nimi)
}

func santsi(update tgbotapi.Update) string {

	kuvaus, err := dbUusinKuppi(update.Message.From.ID)

	if err != nil {
		return err.Error()
	}

	// Lisää uusi kupillinen vanhalla kuvauksella
	kirjaus(update, kuvaus)
	log.Println("Santsi:", update.Message.From.UserName, kuvaus)
	return ""
}

func kupit(id int) string {

	kuppeja := dbKupit(id)

	var txt string
	if kuppeja == 0 {
		txt = "Et taida olla kovin sivistynyt ihminen."
	} else if kuppeja == 1 {
		txt = "Olet juonut yhden kupin kahvia. Tsemppaatko kiitos."
	} else {
		txt = fmt.Sprint("Olet juonut ", kuppeja, " kuppia kahvia.")
	}

	return txt
}

func eiku(update tgbotapi.Update) {

	if update.Message.ReplyToMessage == nil {
		// Tsekkaa, että on vastattu edes johonkin
		log.Println("^ Kusi, ei oo vastattu mihinkään")
		return
	}

	var komento string = update.Message.ReplyToMessage.Command()

	if update.Message.From.ID != update.Message.ReplyToMessage.From.ID {
		// User id:t ei mätsää
		log.Println("^ Kusi, ei vastattu omaan viestiin")
		return
	} else if komento != "kahvi" && komento != "santsi" {
		// Tsekkaa, että vastattu viesti on kirjauskomento
		log.Println("^ Kusi, ei oo kirjauskomento")
		return
	}

	var uusiKuvaus string = update.Message.CommandArguments()
	var userID int = update.Message.From.ID
	var aika int = update.Message.ReplyToMessage.Date

	dbEiku(uusiKuvaus, userID, aika)
}

func poista(update tgbotapi.Update) {

	if update.Message.ReplyToMessage == nil {
		// Tsekkaa, että on vastattu edes johonkin
		log.Println("^ Kusi, ei oo vastattu mihinkään")
		return
	}

	var komento string = update.Message.ReplyToMessage.Command()

	if update.Message.From.ID != update.Message.ReplyToMessage.From.ID {
		// User id:t ei mätsää
		log.Println("^ Kusi, ei vastattu omaan viestiin")
		return
	}
	if komento != "kahvi" && komento != "santsi" {
		// Tsekkaa, että vastattu viesti on kirjauskomento
		log.Println("^ Kusi, ei oo kirjauskomento")
		return
	}

	var userID int = update.Message.From.ID
	var aika int = update.Message.ReplyToMessage.Date

	dbPoista(userID, aika)
}

func viimeisimmat(userID int) string {
	var tekstirivit string = ""
	kirjaukset := dbViimeisimmat(userID, 0)

	for _, s := range kirjaukset {
		tekstirivit += s.teksti + "\n"
	}

	if tekstirivit == "" {
		return "Tyhjältä näyttää.."
	} else {
		return "Sun viimeisimmät kupit:\n" + tekstirivit
	}
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

		if update.CallbackQuery != nil {
			// fmt.Printf("%+v\n", update.CallbackQuery.Message.Chat)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			msg := painallus(update)
			bot.Send(msg)
			continue
		}

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
					"/kupit\nKertoo montako kuppia oot juonu.\n" +
					"/viimeisimmat\nNäyttää max 5 viimeisintä kirjattua kuppia.\n" +
					"/eiku [valinnainen kuvaus]\nVastaamalla aiempaan kirjaukseen vaihtaa sen kuvausta.\n" +
					"/poista\nVastaamalla aiempaan kirjaukseen poistaa sen."
				msgHelp := tgbotapi.NewMessage(chatID, txt)
				bot.Send(msgHelp)
				log.Println("Help:", userName)

			case "kahvi":
				var kuvaus string = update.Message.CommandArguments()

				if len(kuvaus) > 255 {
					log.Println("--Liian pitkä kirjaus:", userName, kuvaus)
					bot.Send(tgbotapi.NewMessage(chatID, "Liian pitkä kuvaus.."))
				} else {
					log.Println("Kirjaus:", userName, kuvaus)
					kirjaus(update, "")
				}

			case "santsi":
				log.Println("Santsi:", userName)
				var err string = santsi(update)
				if err != "" {
					log.Println("Santsi kusee, ei edellistä kuppia", userName)
					errMsg := tgbotapi.NewMessage(chatID, err)
					bot.Send(errMsg)
				}

			case "kupit":
				log.Println("Kupit:", userName)
				txt := kupit(userID)

				msgStats := tgbotapi.NewMessage(chatID, txt)
				bot.Send(msgStats)

			case "eiku":
				log.Println("Eiku:", userName)
				eiku(update)

			case "viimeisimmat":
				log.Println("Viimeisimmät:", userName)
				viesti := viimeisimmat(userID)
				bot.Send(tgbotapi.NewMessage(chatID, viesti))

			case "poista":
				log.Println("Poista:", userName)
				poista(update)

			case "napit":
				log.Println("Napit:", userName)
				viesti := napit(update)
				bot.Send(viesti)
			}
		}
	}
}
