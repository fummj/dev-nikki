package authN

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/models"
)

func TestSetJWTCookie(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ec := app.NewContext(req, rec)

	tx := models.DBC.DB.Begin()
	user := setUpUserData(t, tx)
	ts, err := GenerateJWT(user)
	if err != nil {
		t.Error(err)
	}

	// response-headerにトークンを付与する場合のテスト
	SetJWTCookie(ec, ts)
	header := ec.Response().Header().Get("Set-Cookie")

	extractedCookie, err := http.ParseSetCookie(header)
	if extractedCookie.Value != ts {
		t.Error("生成したaccess-tokenとresponse-headerに設定したaccess-tokenが異なる。")
	}

	// request-headerにトークンを付与する場合のテスト
	SetJWTCookie(req, ts)
	ac, err := ec.Cookie("access_token")
	if err != nil {
		t.Error(err)
	}

	if ac.Value != ts {
		t.Error("生成したaccess-tokenとrequest-headerに設定したaccess-tokenが異なる。")
	}

	tx.Rollback()
}

func TestParseJWTCookie(t *testing.T) {

}
