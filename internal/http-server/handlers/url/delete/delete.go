package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type UrlDeleter interface {
	DeleteURL(alias string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlDeleter
func New(log *slog.Logger, deleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Настраиваем аттрибут для логгера
		const op = "handlers.url.delete.New"

		// Настраиваем логгер с атрибутом
		log = log.With(
			slog.String("op", op),
			// Передаем контекст запроса в метод для получения id запроса
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			// Логируем, если alias пустой
			log.Info("alias is empty")

			// Передаем наверх JSON с описанием ошибки
			render.JSON(w, r, resp.Error("invalid request"))
		}

		err := deleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("deleted url", slog.String("alias", alias))
	}
}
