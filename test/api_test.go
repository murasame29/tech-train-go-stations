package test

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

const (
	testUserID   = "test"
	testPassword = "test"
)

func TestAPI(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error("dbPathのセットに失敗しました。", err)
		return
	}

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Error("DBの作成に失敗しました。", err)
		return
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("DBのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})
	// t.setEnvを使うとよい
	t.Setenv("BASIC_AUTH_USER_ID", testUserID)
	t.Setenv("BASIC_AUTH_PASSWORD", testPassword)

	// if err := os.Setenv("BASIC_AUTH_USER_ID", testUserID); err != nil {
	// 	t.Errorf("テスト用の環境変数のセットに失敗しました: %v", err)
	// 	return
	// }

	// if err := os.Setenv("BASIC_AUTH_PASSWORD", testPassword); err != nil {
	// 	t.Errorf("テスト用の環境変数のセットに失敗しました: %v", err)
	// 	return
	// }

	env, err := env.GetEnv()
	if err != nil {
		t.Error("環境変数の読み込みに失敗しました。", err)
		return
	}

	r := router.NewRouter(todoDB, env)
	srv := httptest.NewServer(r.Mux)
	defer srv.Close()

	testCases := []struct {
		name          string
		setAuthToken  func(*http.Request)
		createRequest func() (*http.Request, error)
		checkResponse func(*testing.T, *http.Response) error
	}{
		{
			name: "healthz",
			createRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, srv.URL+"/healthz", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) error {
				want := "{\"message\":\"OK\"}\n"
				got, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				if string(got) != want {
					return fmt.Errorf("レスポンスの内容が正しくありません。want: %s, got: %s", want, string(got))
				}
				return nil
			},
		},
		{
			name: "todo_success",
			setAuthToken: func(req *http.Request) {
				req.Header.Set("Authorization", "Basic "+base64NewEncoder(testUserID+":"+testPassword))
			},
			createRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, srv.URL+"/todos", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) error {
				want := "{\"todos\":[]}\n"
				got, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				if string(got) != want {
					return fmt.Errorf("レスポンスの内容が正しくありません。want: %s, got: %s", want, string(got))
				}
				return nil
			},
		},
		{
			name: "todo_fail_InvalidToken",
			setAuthToken: func(req *http.Request) {
				req.Header.Set("Authorization", "Basic "+base64NewEncoder(testUserID+":"))
			},
			createRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, srv.URL+"/todos", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) error {
				want := "{\"error\":\"Wrong userId or password\"}\n"
				got, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				if string(got) != want {
					return fmt.Errorf("レスポンスの内容が正しくありません。want: %s, got: %s", want, string(got))
				}
				return nil
			},
		},
		{
			name:         "todo_fail_EmptyToken",
			setAuthToken: func(req *http.Request) {},
			createRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, srv.URL+"/todos", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) error {
				want := "{\"error\":\"UnAuthorized\"}\n"
				got, err := io.ReadAll(resp.Body)

				if err != nil {
					return err
				}

				if string(got) != want {
					return fmt.Errorf("レスポンスの内容が正しくありません。want: %s got: %s", want, string(got))
				}
				return nil
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.createRequest()
			if err != nil {
				t.Errorf("リクエストの作成に失敗しました: %v", err)
				return
			}

			if tc.setAuthToken != nil {
				tc.setAuthToken(req)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの実行に失敗しました: %v", err)
				return
			}
			defer resp.Body.Close()

			if err := tc.checkResponse(t, resp); err != nil {
				t.Errorf("レスポンスのチェックに失敗しました: %v", err)
				return
			}
		})
	}
}

func base64NewEncoder(rawToken string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawToken))
}
