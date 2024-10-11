package tgalert

import (
	"bytes"
	"encoding/json"
	"monitor/config"
	"net/http"
)

type RemoteAlert struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func SendAlert(msg string) error {
	var Config config.Conf
	Config.GetConfig()

	v := RemoteAlert{
		ChatID: Config.TgBot.ChatID,
		Text:   msg,
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://api.telegram.org/bot"+Config.TgBot.Token+"/sendMessage", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
