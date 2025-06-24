package login

import (
	"testing"

	"gorm.io/gorm"

	"dev_nikki/internal/models"
)

var (
	test_username string = "test_login"
	test_email    string = "test@login.com"
	test_password string = "398asdlkfj$Kad"
	test_salt     string = "aeijae3$d.I9"
)

func setUpUserData(t *testing.T, tx *gorm.DB) *models.User {
	_, user, err := models.CreateUser(tx, test_username, test_email, test_password, test_salt)
	if err != nil {
		t.Error(err)
	}
	return user
}

func Test_verifyHashedPassword(t *testing.T) {
	tx := models.DBC.DB.Begin()
	user := setUpUserData(t, tx)
	if err := verifyHashedPassword(test_password, user); err != nil {
		t.Error(err)
	}
}
