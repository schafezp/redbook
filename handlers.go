package redbook

import (
	"database/sql"
	db "database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Handler main point of entry for handlers
func Handler(bot tgbotapi.BotAPI, update tgbotapi.Update, db *db.DB) {

	print("Handling bot api")
	log.Printf("Is user making prediction %t", isUserMakingPrediction(update.Message.From.ID, db))

	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			msg.Text = "Starting message"
			bot.Send(msg)
		case "predict":
			startPrediction(bot, update, db)
		default:
			msg.Text = "Unknown command"
			bot.Send(msg)
		}
	} else {
		if isUserMakingPrediction(update.Message.From.ID, db) {
			submitPrediction(*update.Message, db, bot)
		}
	}

}

//TODO: add command to undo

//passed in
func startPrediction(bot tgbotapi.BotAPI, update tgbotapi.Update, db *db.DB) {
	msgText := "State prediction include date in YYYY-MM-DD form"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
	db.Query("INSERT INTO prediction_in_progress VALUES ($1)", update.Message.From.ID)
}

func isUserMakingPrediction(userID int, db *db.DB) bool {
	rows := db.QueryRow("SELECT EXISTS(SELECT userid FROM prediction_in_progress WHERE userid = $1)", userID)
	var isUserMakingPred bool
	switch err := rows.Scan(&isUserMakingPred); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return true
	case nil:
		fmt.Println("User pred ", isUserMakingPred)
		return isUserMakingPred
	default:
		panic(err)
	}

}

func submitPrediction(message tgbotapi.Message, db *db.DB, bot tgbotapi.BotAPI) {

	db.Query("DELETE FROM prediction_in_progress WHERE userid=$1", message.From.ID)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Submitted prediction")
	bot.Send(msg)

}
