package delete

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

type IDelete interface {
	DeleteUrl(alias string) (string, error)
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=IRedirect

func New(log *slog.Logger, deleteUrlGetter IDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete"

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

		// execute func DeleteUrl from sqlite storage
		resAlias, err := deleteUrlGetter.DeleteUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Error("Url not found", sl.Err(err))
			render.JSON(w, r, resp.Error("Url not found"))
			return
		}

		if err != nil {
			log.Error("Failed to delete url by alias", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("url deleted", slog.String("alias", resAlias))

		responseOk(w, r, resAlias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:       alias,
	})
}
