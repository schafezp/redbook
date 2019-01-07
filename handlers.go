package redbook

import (
	"database/sql"
	db "database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//TODO: add command to undo

//passed in
func startPrediction(bot tgbotapi.BotAPI, update tgbotapi.Update, db *db.DB) {
	msgText := "State prediction include date in YYYY-MM-DD form"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
	db.Query("INSERT INTO prediction_in_progress VALUES ($1)", update.Message.From.ID)
}

func isUserMakingPrediction(userID int, db *db.DB) bool {
	pred := CreatingPrediction{}
	row, _ := db.Query("SELECT userid FROM prediction_in_progress WHERE userid = $1", userID)
	err := row.Scan(&pred.UserId)
	//if err there are no rows, user not making prediction
	return err != sql.ErrNoRows
}

func submitPrediction(message tgbotapi.Message, db *db.DB, bot tgbotapi.BotAPI) {

	db.Query("DELETE prediction_in_progress WHERE userid=$1", message.From.ID)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Submitted prediction")
	bot.Send(msg)

}

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

			/* for rows.Next() {
				prediction := redbook.CreatingPrediction{}
				rows.Scan(&prediction.UserId)
				println("userid: " + strconv.Itoa(prediction.UserId))
			} */

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
