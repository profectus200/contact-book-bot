package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/profectus200/contact-book-bot/internal/model/callbacks"
)

var editContactKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change name", callbacks.ChangeContactName),
		tgbotapi.NewInlineKeyboardButtonData("Change phone", callbacks.ChangeContactPhone),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change birthday", callbacks.ChangeContactBirthday),
		tgbotapi.NewInlineKeyboardButtonData("Change description", callbacks.ChangeContactDescription),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Delete contact", callbacks.DeleteContact),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Save", callbacks.ChangeContactDone),
	),
)
