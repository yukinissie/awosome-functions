package line

import (
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func Excute(message string) error {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return err
	}

	messages := []linebot.SendingMessage{linebot.NewTextMessage(message)}

	_, err = bot.BroadcastMessage(messages...).Do()
	if err != nil {
		return err
	}

	return nil
}
