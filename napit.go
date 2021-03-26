package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var defaultKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Santsi", "santsi"),
		tgbotapi.NewInlineKeyboardButtonData("Tilastot", "tilastot"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Poista", "poista"),
	),
)

var backKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Takaisin", "aloitus"),
	),
)

var defaultText = "Hyvää päivää. Mitä sais olla?"

func generoiSantsiNapit(nimet []string) tgbotapi.InlineKeyboardMarkup {
	var napit [][]tgbotapi.InlineKeyboardButton
	for _, teksti := range nimet {
		var data string = "santsi " + teksti
		napit = append(napit,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(teksti, data)))
	}

	napit = append(napit, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Takaisin", "aloitus")))

	return tgbotapi.NewInlineKeyboardMarkup(napit...)
}

func generoiPoistaNapit(kahvit []string) tgbotapi.InlineKeyboardMarkup {
	var napit [][]tgbotapi.InlineKeyboardButton
	for _, teksti := range kahvit {

		timestampInt := strings.SplitN(teksti, " ", 2)[0]
		kahviTeksti := strings.SplitN(teksti, " ", 2)[1]

		timestamp, err := strconv.ParseInt(timestampInt, 10, 64)
		if err != nil {
			panic(err)
		}
		var aika time.Time = time.Unix(timestamp, 0)
		loc, _ := time.LoadLocation("Europe/Helsinki")
		aikaTeksti := aika.In(loc).Format("2.1. 15:04")
		var nappiTeksti string = "Poista " + kahviTeksti + " (" + aikaTeksti + ")"
		var nappiData string = "poista " + fmt.Sprint(timestamp) + " " + kahviTeksti
		napit = append(napit,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(nappiTeksti, nappiData)))
	}

	napit = append(napit, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Takaisin", "aloitus")))

	return tgbotapi.NewInlineKeyboardMarkup(napit...)
}

func napit(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	var msg tgbotapi.MessageConfig

	if update.Message.Chat.Type == "private" {
		msg = tgbotapi.NewMessage(chatID, defaultText)
		msg.ReplyMarkup = defaultKeyboard
	} else {
		msg = tgbotapi.NewMessage(chatID, "Privaan...")
	}
	return msg
}

func painallus(update tgbotapi.Update) tgbotapi.EditMessageTextConfig {
	if update.CallbackQuery.Message.Chat.Type != "private" {
		log.Println("Callback not in private")
		os.Exit(1)
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	userID := int(chatID)
	msgID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data
	var kb tgbotapi.InlineKeyboardMarkup
	var text string

	log.Println(userID, data)

	switch data {
	case "aloitus":
		kb = defaultKeyboard
		text = defaultText

	case "santsi":
		kahvit := dbViimeisimmatUniikit(userID)
		kb = generoiSantsiNapit(kahvit)
		text = "Mitäs laitetaan?"

	case "tilastot":
		kb = backKeyboard
		text = fmt.Sprint("Olet juonut ", dbKupit(userID), " kuppia kahvia. ")

	case "poista":
		kahvit := dbViimeisimmat(userID, true)
		kb = generoiPoistaNapit(kahvit)
		text = "Minkä kahvin haluat unohtaa?"

	default:

		if strings.Contains(data, "santsi ") {
			kahvi := data[7:]
			dbKirjaus(userID, int(time.Now().Unix()), kahvi, "")
			kb = defaultKeyboard
			text = "Santsattu, mitäs sitte?"

		} else if strings.Contains(data, "poista ") {
			// dbKirjaus(userID, int(time.Now().Unix()), kahvi, "")
			aikaStr := strings.SplitN(data, " ", 3)[1]
			aika, err := strconv.Atoi(aikaStr)
			if err != nil {
				panic(err)
			}
			dbPoista(userID, aika)
			kb = defaultKeyboard
			text = "Poistettu, mitäs sitte?"
		}
	}

	editKb := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, kb)

	return editKb
}
