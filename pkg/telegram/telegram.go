package telegram

import (
	"amplify_bot/pkg/ffmpeg"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(bot *tgbotapi.BotAPI) *Telegram {
	return &Telegram{
		bot: bot,
	}
}

func (t Telegram) Process(update *tgbotapi.Update) {

	if update.Message.Voice != nil {

		file, err := t.bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Voice.FileID})
		if err != nil {
			return
		}

		downloadedVoiceMessage, err := ioutil.TempFile("", "vm.*.oga")

		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.bot.Token, file.FilePath)

		err = downloadFile(downloadedVoiceMessage.Name(), url)
		if err != nil {
			return
		}

		defer os.Remove(downloadedVoiceMessage.Name())

		amplified, err := ffmpeg.Amplify(downloadedVoiceMessage.Name())
		if err != nil {
			log.Printf("error in amplify: %v", err)
			return
		}

		defer os.Remove(amplified.Name())

		responseBytes, err := ioutil.ReadFile(amplified.Name())
		if err != nil {
			fmt.Println("error in read: %v", err)
			return
		}

		msg := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "amplified.oga", Bytes: responseBytes})
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = t.bot.Send(msg)
		if err != nil {
			log.Printf("error in send: %v", err)
			return
		}
	}

}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
