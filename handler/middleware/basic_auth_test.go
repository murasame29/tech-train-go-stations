package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/env"
)

// basic_auth middlewarer test
func TestBasicAuth(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		authToken      string
		expectedStatus int
	}{
		{
			name:           "authTokenが空の場合",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenがbase64でdecodeできない場合",
			authToken:      "test",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenがbase64でdecodeできる場合",
			authToken:      "dGVzdDp0ZXN0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "authTokenがbase64でdecodeできる場合",
			authToken:      "dGVzdDp0ZXN0",
			expectedStatus: http.StatusOK,
		},
	}
	// テストケースを回す
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テストケースの設定
			req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tc.authToken)
			// テスト対象の関数を呼び出す
			rr := httptest.NewRecorder()
			m := NewMiddleware(&env.Env{UserID: "test", Password: "test"})
			m.BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(rr, req)
			// テスト結果の評価
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected: %d, actual: %d", tc.expectedStatus, rr.Code)
			}
		})
	}

}
