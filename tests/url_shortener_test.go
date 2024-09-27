package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"usekit-go/internal/http-server/handlers/url/save"
	"usekit-go/internal/lib/api"
	"usekit-go/internal/lib/random"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	// создается базовый урл, по которому клиент будет обращаться
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	// создается клиент, через который будут отправляться запросы
	e := httpexpect.Default(t, u.String())

	// отправляется post запрос из созданного клиента
	// передаем объект запроса из хендлера save, который маршалится в json
	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("user", "password").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save
			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("user", "password").
				Expect().Status(http.StatusOK).
				JSON().Object()

			// если в тест-кейсе ожидается ошибка, то ожидаем, что в ответе не будет alias
			if tc.error != "" {
				resp.NotContainsKey("alias")

				// проверяем соответствие полученной ошибке с ошибкой из тест-кейса
				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			// если в тест-кейсе присутствовал alias, то проверяем соответсвие его к полученному
			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				// иначе проверяем наличие рандомно сгенерированного alias
				resp.Value("alias").String().NotEmpty()

				// записываем сгенерированный alias локальной переменной alias
				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			testRedirect(t, alias, tc.url)

			// Delete
			deletedAlias := resp.Value("alias").String().Raw()

			respDelete := e.DELETE("/url/"+deletedAlias).
				WithBasicAuth("user", "password").
				Expect().Status(http.StatusOK).
				JSON().Object()

			respDelete.Value("status").String().IsEqual("OK")

			// Redirect again
			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	uString := u.String()
	redirectedToURL, err := api.GetRedirect(uString)
	require.NoError(t, err)

	// проверяем соответствие сгенерированного url к url, полученному в параметре
	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}