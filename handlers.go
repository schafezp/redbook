package redbook

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handler(bot tgbotapi.BotAPI, update tgbotapi.Update) {
	print("Handling bot api")
	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			msg.Text = "Starting message"
			bot.Send(msg)
		default:
			msg.Text = "Unknown command"
			bot.Send(msg)
		}
	}
}
