package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var defaultKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Santsi", "santsi"),
		tgbotapi.NewInlineKeyboardButtonData("Tilastot", "tilastot"),
	),
)

var backKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Takaisin", "aloitus"),
	),
)

var defaultText = "Hyvää päivää. Mitä sais olla?"

func generoiSantsiNapit(nimet []string, datat []string) tgbotapi.InlineKeyboardMarkup {
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

	switch data {
	case "aloitus":
		kb = defaultKeyboard
		text = defaultText

	case "santsi":
		kahvit := dbViimeisimmatUniikit(userID)
		kb = generoiSantsiNapit(kahvit, kahvit)
		text = fmt.Sprint("Mitäs laitetaan?")

	case "tilastot":
		kb = backKeyboard
		text = fmt.Sprint("Olet juonut ", dbKupit(userID), " kuppia kahvia. ")

	default:
		if strings.Contains(data, "santsi ") {
			kahvi := data[7:]
			log.Println(kahvi)
			dbKirjaus(userID, int(time.Now().Unix()), kahvi, "")
			kb = defaultKeyboard
			text = "Santsattu, mitäs sitte?"
		}
	}

	editKb := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, kb)

	return editKb
}
