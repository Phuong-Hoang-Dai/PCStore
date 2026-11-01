package model_test

import (
	"testing"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

func TestUser(t *testing.T) {
	var user model.User
	password := []byte("Test12345**")
	user.SetHashPassword(password)

	t.Run("correct password verifies", func(t *testing.T) {
		if err := user.VerifyPassword(password); err != nil {
			t.Errorf("VerifyPassword: %v", err)
		}
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		if err := user.VerifyPassword([]byte("wrong")); err == nil {
			t.Error("VerifyPassword accepted a wrong password")
		}
	})
}
