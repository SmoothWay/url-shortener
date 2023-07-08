package save

import (
	"net/http"

	resp "github.com/SmoothWay/url-shortener/internal/lib/api/response"
	"github.com/SmoothWay/url-shortener/internal/lib/logger/sl"
	"github.com/SmoothWay/url-shortener/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))

			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
	}
}
