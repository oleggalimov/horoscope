package radioc

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"horoscope/internal/horoscope/model"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	URL = "https://www.radioc.ru/programs/357/"
)

type RadioCiProvider struct{}

func NewRadioCiProvider() *RadioCiProvider {
	return &RadioCiProvider{}
}

func (*RadioCiProvider) GetDailyHoroscope() (model.DailyHoroscope, error) {
	httpResponse, err := getHtmlPage()
	if err != nil {
		return model.DailyHoroscope{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Ошибка при обращении к ресурсам Радио Си", err)
		}
	}(httpResponse.Body)

	return parseHoroscopePage(httpResponse)
}

func getHtmlPage() (*http.Response, error) {
	resp, err := http.Get(URL)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении гороскопа Радио Си: %w", err)
	}
	return resp, nil
}

func parseHoroscopePage(response *http.Response) (model.DailyHoroscope, error) {

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return model.DailyHoroscope{}, fmt.Errorf("ошибка парсинга: %w", err)
	}

	forecasts := make([]model.ZodiacForecast, 0)
	doc.Find(".astro_main").Each(
		func(i int, selection *goquery.Selection) {
			fc, err := mapForeCast(selection)
			if err == nil {
				forecasts = append(forecasts, fc)
			}
		})

	return model.DailyHoroscope{
		Source:    "Radio C",
		Date:      time.Now(),
		Forecasts: forecasts,
	}, nil
}

func mapForeCast(s *goquery.Selection) (model.ZodiacForecast, error) {
	text := strings.TrimSpace(s.Text())
	if text == "" {
		return model.ZodiacForecast{}, errors.New("Пустой astro_main block")
	}

	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		return model.ZodiacForecast{}, fmt.Errorf("Неверный astro_main format: %q", text)
	}

	zodiac := strings.TrimSpace(lines[0])
	forecast := strings.TrimSpace(strings.Join(lines[1:], " "))

	if zodiac == "" || forecast == "" {
		return model.ZodiacForecast{}, errors.New("Знак или прогноз - пустой")
	}

	return model.ZodiacForecast{
		Zodiac: model.ZodiacSign(zodiac),
		Text:   forecast,
	}, nil
}
