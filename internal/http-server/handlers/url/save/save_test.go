package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"usekit-go/internal/http-server/handlers/url/save"
	"usekit-go/internal/http-server/handlers/url/save/mocks"
	"usekit-go/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// пример запроса в формате json
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// создаем новый запрос с телом типом метода, урлом и телом запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			// сгенерируется и сразу выбросится ошибка при неудачном создании запроса
			require.NoError(t, err)

			// response recorder, в который записывается ответ нашего хендлера
			rr := httptest.NewRecorder()
			// выполняем запрос
			handler.ServeHTTP(rr, req)

			// проверяем соответствие статус кодов успешного ответа
			require.Equal(t, rr.Code, http.StatusOK)

			// получаем тело ответа
			body := rr.Body.String()

			var resp save.Response

			// проверяем отсутствие ошибок при парсинге json body ответа
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// проверяем соответствие ошибки, которую вернул хендлер ошибке, которая определна в текущем тест кейсе
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}