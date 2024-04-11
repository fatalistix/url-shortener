package del

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/fatalistix/url-shortener/internal/lib/api/response"
	"github.com/fatalistix/url-shortener/internal/lib/logger/sl"
	"github.com/fatalistix/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
}

// URLDeleter is an interface for deleting url by alias
//
//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)

				render.JSON(w, r, response.Error("not found"))

				return
			} else {
				log.Error("failed to delete url", sl.Err(err))

				render.JSON(w, r, response.Error("failed to delete url"))

				return
			}
		}

		log.Info("url deleted")

		render.JSON(w, r, Response{
			response.OK(),
		})
	}
}
