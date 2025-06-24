package signup

import (
	"testing"
	"unicode/utf8"
)

func TestGetPepper(t *testing.T) {
	s := GetPepper()
	if s == "" {
		t.Error("pepperを取得できていない。")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt := GenerateSalt()
	if salt == "" {
		t.Error("saltが生成できていない。")
	}
}

func TestPasswordHashing(t *testing.T) {
	salt := GenerateSalt()
	data := []struct {
		name     string
		password string
	}{
		{"success", "testtE4$"},
		{"failed", "testte4$"},
		{"failed", "testtEs$"},
		{"failed", "testtE4s"},
		{"failed", "testtE4"},
		{"failed", "testT4$;asdfdfadsfadfasdfadfdfadsfadsfadfadsfadfadfadfadfadfadfadfadfadfadfadfadfadfadfadsfadfadfsa"},
		{"failed", "testT4$;asdfdfadsfadfasdfadfdfadsfadsfadfadsfadfadfadfadfadfadfadfadfadfadfadfadfadfadfadsfadfadfsas"},
	}

	sc := utf8.RuneCountInString(GenerateSalt())
	gc := utf8.RuneCountInString(GetPepper())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {

			totalCount := sc + gc + utf8.RuneCountInString(d.password)
			hashed, err := PasswordHashing(salt, d.password)
			hc := utf8.RuneCountInString(hashed)
			if d.name == "success" && err != nil {
				t.Error(err)
			}

			// もしsalt, password, pepperの文字数を足しただけの文字数ならエラー扱い
			if d.name == "failed" && hc == totalCount {
				t.Error("パスワードハッシュ化に失敗しました")
			}
		})
	}
}
