package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/fatalistix/url-shortener/internal/lib/api/response"
	"github.com/fatalistix/url-shortener/internal/lib/logger/sl"
	"github.com/fatalistix/url-shortener/internal/lib/random"
	"github.com/fatalistix/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(aliasLength int, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// Если запрос пришел с пустым телом
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.JSON(w, r, response.Error("empty request"))

				return
			} else {
				log.Error("failed to decode request body", sl.Err(err))

				render.JSON(w, r, response.Error("failed to decode request"))

				return
			}
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		// TODO: check that links are unique
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))

				render.JSON(w, r, response.Error("url already exists"))

				return
			} else {
				log.Error("failed to add url", sl.Err(err))

				render.JSON(w, r, response.Error("failed to add url"))

				return
			}
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
