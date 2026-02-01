package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"horoscope/internal/telegram/model"
	"io"
	"log"
	"net/http"
)

type bot struct {
	token      string
	updatesUrl string
}

func newBot(token string) *bot {
	updatesUrl := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=%d", token, 0)
	return &bot{token, updatesUrl}
}

func (bot *bot) getUpdates() (*[]model.Update, error) {

	var response, err = http.Get(bot.updatesUrl)
	if err != nil {
		log.Printf("Не удалось получить апдейт от сервера: %v", err)
		return nil, err
	}

	bodyBytes, _ := io.ReadAll(response.Body)
	response.Body.Close()
	fmt.Printf("Ответ от сервера: %s\n", string(bodyBytes))

	updates := model.UpdatesResponse{}
	unmarshallErr := json.Unmarshal(bodyBytes, &updates)
	if unmarshallErr != nil {
		log.Printf("Ошибка парсинга: %v\nОтвет: %s", unmarshallErr, string(bodyBytes))
		return nil, err
	}

	if !updates.IsOk {
		log.Println("Апи telegram вернул ответ с флагом не Ок")
		return nil, errors.New("неуспешный ответ telegram")
	}

	return &updates.Result, nil
}
