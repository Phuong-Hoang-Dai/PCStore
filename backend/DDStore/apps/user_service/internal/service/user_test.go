package service_test

import (
	"testing"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name    string
		uCD     service.UserCreateDTO
		uD      service.UserDTO
		wantErr bool
	}{
		{
			name:    "test happy case",
			uCD:     service.UserCreateDTO{Name: "HoangDai", Password: "Dai12345678"},
			uD:      service.UserDTO{Name: "HoangDai", Password: "Dai12345678"},
			wantErr: false,
		},
		{
			name:    "test wrong password",
			uCD:     service.UserCreateDTO{Name: "HoangDai", Password: "Dai12345678"},
			uD:      service.UserDTO{Name: "HoangDai", Password: "ThisIsWrongPassword"},
			wantErr: true,
		},
	}

	var err error
	var mock service.MockRepos
	mock.Init()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.uD.Id, err = service.CreateUser(test.uCD, mock)
			assert.NoError(t, err)
			_, err := service.VerifyPassword(test.uD, mock)
			if test.wantErr {
				assert.ErrorIs(t, err, model.ErrUserNameOrPasswordIncorrect)

			} else {
				assert.NoError(t, err)
			}
		})
	}
}
