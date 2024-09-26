package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "usekit-go/internal/lib/api/response"
	"usekit-go/internal/lib/logger/sl"
	"usekit-go/internal/storage"
)

type IRedirect interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter IRedirect) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect"

		// mark next logs by request id
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var alias = chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("Missing parameter 'alias'")
			render.JSON(w, r, resp.Error("Missing parameter 'alias'"))
			return
		}
		log.Info("alias received")

		// execute func GetUrl from sqlite storage
		resUrl, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Error("Url not found", err)
			render.JSON(w, r, resp.Error("Url not found"))
			return
		}

		if err != nil {
			log.Error("Failed to get url by alias", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("got url", slog.String("url", resUrl))

		// redirect to found url
		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
