package delete

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	_ "url-shortener/internal/lib/api/response"
	_ "url-shortener/internal/lib/logger/sl"
)

func TestNew(t *testing.T) {
	// Создаем мок UrlDeleter

	tests := []struct {
		name          string
		alias         string
		expectedError error
		expectedCode  int
	}{
		{
			name:          "Successful deletion",
			alias:         "test_alias",
			expectedError: nil,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Deletion error",
			alias:         "test_alias",
			expectedError: errors.New("delete error"),
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:          "Empty alias",
			alias:         "",
			expectedError: nil,
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Очищаем вызовы в моке перед каждым тестом
			urlDeleteMock := mocks.NewUrlDeleter(t)

			// Создаем новый роутер Chi
			r := chi.NewRouter()

			// Устанавливаем обработчик для маршрута
			r.Delete("/delete/{alias}", New(nil, urlDeleteMock))

			// Создаем тестовый запрос DELETE к маршруту с указанным псевдонимом
			req, _ := http.NewRequest("DELETE", "/delete/"+test.alias, nil)

			// Создаем запись для записи ответа сервера
			rec := httptest.NewRecorder()

			// Выполняем запрос
			r.ServeHTTP(rec, req)

			// Проверяем код ответа
			assert.Equal(t, test.expectedCode, rec.Code)

			// TODO: Добавьте дополнительные проверки, в зависимости от требований вашего кода
		})
	}
}
