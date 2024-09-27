package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"usekit-go/internal/http-server/handlers/url/redirect"
	"usekit-go/internal/http-server/handlers/url/redirect/mocks"
	"usekit-go/internal/lib/api"
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
			alias: "google",
			url:   "https://www.google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewIRedirect(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetUrl", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			// получили url, на который произошел редирект
			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			// проверяем отсутствие ошибок при выполнении операции получения урла на редирект
			require.NoError(t, err)

			// проверяем соответствие полученного url тому, который находится в тест-кейсе
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
