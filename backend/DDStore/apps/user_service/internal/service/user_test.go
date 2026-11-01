package service_test

import (
	"context"
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
			test.uD.Id, err = service.CreateUser(test.uCD, mock, context.Background())
			assert.NoError(t, err)
			_, err := service.VerifyPassword(test.uD, mock, context.Background())
			if test.wantErr {
				assert.ErrorIs(t, err, model.ErrUserNameOrPasswordIncorrect)

			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateGuestUser(t *testing.T) {
	var mock service.MockRepos
	mock.Init()
	ctx := context.Background()

	_, err := service.CreateUser(service.UserCreateDTO{Name: "Khách"}, mock, ctx)
	assert.ErrorIs(t, err, model.ErrInvalidEmail, "a guest needs an email")

	_, err = service.CreateUser(service.UserCreateDTO{Email: "guest@example.com"}, mock, ctx)
	assert.ErrorIs(t, err, model.ErrNameRequired)

	first, err := service.CreateUser(service.UserCreateDTO{Name: "Khách", Email: " Guest@Example.com "}, mock, ctx)
	assert.NoError(t, err)
	again, err := service.CreateUser(service.UserCreateDTO{Name: "Khách", Email: "guest@example.com"}, mock, ctx)
	assert.NoError(t, err)
	assert.Equal(t, first, again, "the same email gets the same account")

	_, err = service.VerifyPassword(service.UserDTO{Name: "Khách"}, mock, ctx)
	assert.Error(t, err, "a guest account can't log in")
}
