package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fatalistix/url-shortener/internal/http-server/handlers/url/save"
	"github.com/fatalistix/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/fatalistix/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		mockError   error
		name        string
		alias       string
		url         string
		respError   string
		aliasLength int
	}{
		{
			name:        "Success",
			alias:       "test_alias",
			url:         "https://google.com",
			aliasLength: 10,
		},
		{
			name:        "Empty alias",
			alias:       "",
			url:         "https://google.com",
			aliasLength: 10,
		},
		{
			name:        "Empty URL",
			url:         "",
			alias:       "some_alias",
			respError:   "field URL is a required field",
			aliasLength: 10,
		},
		{
			name:        "SaveURL Error",
			alias:       "test_alias",
			url:         "https://google.com",
			respError:   "failed to add url",
			mockError:   errors.New("unexpected error"),
			aliasLength: 10,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(tc.aliasLength, slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
