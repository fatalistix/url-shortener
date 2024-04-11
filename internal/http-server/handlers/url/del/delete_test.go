package del_test

import (
	"encoding/json"
	"github.com/fatalistix/url-shortener/internal/http-server/handlers/url/del"
	"github.com/fatalistix/url-shortener/internal/http-server/handlers/url/del/mocks"
	"github.com/fatalistix/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
			}

			handler := del.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			request, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, http.NoBody)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Delete("/url/{alias}", handler)

			router.ServeHTTP(recorder, request)

			require.Equal(t, recorder.Code, http.StatusOK)

			body := recorder.Body.String()

			var resp del.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
