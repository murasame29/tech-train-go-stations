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
			authToken:      "Basic ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenがbase64でdecodeできない場合",
			authToken:      "Basic test",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenにuserIDのみ設定されている場合",
			authToken:      "Basic dGVzdDo=",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenにpasswordのみ設定されている場合",
			authToken:      "Basic OnRlc3Q=",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "authTokenがbase64でdecodeできる場合",
			authToken:      "Basic dGVzdDp0ZXN0",
			expectedStatus: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tc.authToken)

			rr := httptest.NewRecorder()
			m := NewMiddleware(&env.Env{UserID: "test", Password: "test"})
			m.BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected: %d, actual: %d", tc.expectedStatus, rr.Code)
			}
		})
	}

}
