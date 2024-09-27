package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	// создаем кастомного клиента и проверяем направление редиректов без выполнения самого процесса редиректа
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after 1st redirect
		},
	}

	// делаем запрос на переданный параметром url
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	// закрываем тело ответа
	defer func() { _ = resp.Body.Close() }()

	// проверяем совпадение статус-кода ответа с ожидаемым
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	// возвращаем url, на который происходит редирект
	// берем его из хедера Location ответа
	return resp.Header.Get("Location"), nil
}