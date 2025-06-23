package authN

import (
	"testing"
)

var (
	test254Email string = "test254@adsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsasdfasdfasdfas.com"
	test255Email string = "test255@adsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsadsfaefadsfadsfadsfadfadfadfadfadfadfadfadfadfadfadfadfdsasdfasdfasdfasd.com"
	testPassword string = ""
)

func TestEmailValidation(t *testing.T) {

	data := []struct {
		name  string
		email string
	}{
		{"success", test_email},
		{"success", "test@test.co"},
		{"failed", "testtest.com"},
		{"failed", "test@test"},
		{"failed", "test@test.c"},
		{"success", test254Email},
		{"failed", test255Email},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := EmailValidation(d.email)
			if d.name == "success" && err != nil {
				t.Error(err)
			}

			if d.name == "false" && err == nil {
				t.Error("エラーにならないといけない箇所なのに処理が通っている。")
			}
		})
	}
}

func TestPasswordValidation(t *testing.T) {

	data := []struct {
		name     string
		password string
	}{
		{"success", "testT4$;"},
		{"failed", "tesT$4;"},
		{"failed", "testTe$;"},
		{"failed", "testTest"},
		{"failed", "testT4$;asdfdfadsfadfasdfadfdfadsfadsfadfadsfadfadfadfadfadfadfadfadfadfadfadfadfadfadfadsfadfadfsa"},
		{"failed", "testT4$;asdfdfadsfadfasdfadfdfadsfadsfadfadsfadfadfadfadfadfadfadfadfadfadfadfadfadfadfadsfadfadfsas"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := PasswordValidation(d.password)
			if d.name == "success" && err != nil {
				t.Error(err)
			}

			if d.name == "false" && err == nil {
				t.Error("エラーにならないといけない箇所なのに処理が通っている。")
			}
		})
	}
}
