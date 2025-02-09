// learn about telegram bro. and http/https, TCP/UDP and other
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const token = "Api Token"
const apUrl = "https://api.telegram.org/bot" + token

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}
type Chat struct {
	ID int `json:"id"`
}

func GetUpdates(offfset int) ([]Update, error) {
	urll := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=25", apUrl, offfset)
	resp, err := http.Get(urll)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Result []Update `json:"result"`
	}
	json.Unmarshal(body, &result)

	return result.Result, nil
}

type SendMessageRequest struct {
	ChtID int    `json:"chat_id"`
	Text  string `json:"text"`
}

func SendMessage(chatID int, text string) error {
	url := fmt.Sprintf("%s/sendMessage", apUrl)
	msg := SendMessageRequest{ChtID: chatID, Text: text}
	jsonDating, _ := json.Marshal(msg)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonDating))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	offset := 0

	for {
		updates, err := GetUpdates(offset)
		if err != nil {
			fmt.Println("Error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message.Text == "/start" {
				fmt.Println("Получена команда /start, отправляю ответ...")
				SendMessage(update.Message.Chat.ID, "Привет! я пидорас тупой, пососи хуйца")
			}
			offset = update.UpdateID + 1
		}
	}
}
