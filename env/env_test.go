package env

import (
	"os"
	"testing"
)

// GetEnv test
func TestGetEnv(t *testing.T) {
	testCases := []struct {
		name        string
		userId      string
		password    string
		expectedEnv *Env
		expectedErr error
	}{
		{
			name:        "環境変数が正しく設定されている場合",
			userId:      "test",
			password:    "test",
			expectedEnv: &Env{UserID: "test", Password: "test"},
			expectedErr: nil,
		},
		{
			name:        "環境変数が正しく設定されていない場合",
			userId:      "",
			password:    "",
			expectedEnv: nil,
			expectedErr: &EnvError{EnvName: USER_ID},
		},
		{
			name:        "環境変数が正しく設定されていない場合",
			userId:      "test",
			password:    "",
			expectedEnv: nil,
			expectedErr: &EnvError{EnvName: PASSWORD},
		},
		{
			name:        "環境変数が正しく設定されていない場合",
			userId:      "",
			password:    "test",
			expectedEnv: nil,
			expectedErr: &EnvError{EnvName: USER_ID},
		},
	}
	// テストケースを回す
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 環境変数の設定
			os.Setenv(USER_ID, tc.userId)
			os.Setenv(PASSWORD, tc.password)
			// テスト対象の関数を呼び出す
			actualEnv, actualErr := GetEnv()
			// テスト結果の評価
			if actualEnv != nil {
				if actualEnv.UserID != tc.expectedEnv.UserID {
					t.Errorf("expected: %s, actual: %s", tc.expectedEnv.UserID, actualEnv.UserID)
				}
				if actualEnv.Password != tc.expectedEnv.Password {
					t.Errorf("expected: %s, actual: %s", tc.expectedEnv.Password, actualEnv.Password)
				}
			}
			if actualErr != nil {
				if actualErr.Error() != tc.expectedErr.Error() {
					t.Errorf("expected: %s, actual: %s", tc.expectedErr.Error(), actualErr.Error())
				}
			}
		})
	}
}
