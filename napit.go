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

var backButton = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Takaisin", "aloitus"),
)

var backKeyboard = tgbotapi.NewInlineKeyboardMarkup(backButton)

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

func generoiKahviNapit(kirjaukset []Kirjaus, etuliite string) [][]tgbotapi.InlineKeyboardButton {
	var napit [][]tgbotapi.InlineKeyboardButton

	for _, k := range kirjaukset {
		var data = etuliite + " " + fmt.Sprint(k.timestamp) + " " + k.teksti
		napit = append(napit, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(k.teksti, data)))
	}

	return napit
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
		kahvit := dbViimeisimmat(userID, 0)
		napit := generoiKahviNapit(kahvit, "poista")
		napit = append(napit, backButton)
		kb = tgbotapi.NewInlineKeyboardMarkup(napit...)
		text = "Minkä kahvin haluat unohtaa?"

	default:

		if strings.Contains(data, "santsi ") {
			kahvi := data[7:]
			dbKirjaus(userID, int(time.Now().Unix()), kahvi, "")
			kb = defaultKeyboard
			text = "Santsattu, mitäs sitte?"

		} else if strings.Contains(data, "poista ") {
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
