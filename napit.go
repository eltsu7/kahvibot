package main

import (
	"fmt"
	"log"
	"os"

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
		kb = backKeyboard
		text = fmt.Sprint("Hienosti! Olet juonut ", dbKupit(userID), " kuppia kahvia.")

	case "tilastot":
		kb = backKeyboard
		kahvit := dbViim(userID)
		text = fmt.Sprint("Olet juonut ", dbKupit(userID), " kuppia kahvia. ", kahvit[0])
	}

	editKb := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, kb)

	return editKb
}
