package delete

import (
	"errors"
	"net/http"

	resp "github.com/SmoothWay/url-shortener/internal/lib/api/response"
	"github.com/SmoothWay/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// type Response struct {
// 	resp.Response
// }

type URLDeleter interface {
	DeleteURL(alias string) error
}

func Delete(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.Delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {

			log.Error("alias is empty")
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, "invalid request")

			return
		}

		err := deleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("alias not found")
				w.WriteHeader(http.StatusNotFound)

				render.JSON(w, r, "no such alias")

				return
			}

			log.Error("failed to get url:", err)
			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, "internal error")

			return
		}

		log.Info(alias, "deleted")

		render.JSON(w, r, resp.OK())
	}

}
