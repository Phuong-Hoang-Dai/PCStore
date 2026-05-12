package utils_test

import (
	"testing"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/model"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {
	tests := []struct {
		uR model.UseResp
	}{
		{model.UseResp{Id: 999666, Name: "Keyone"}},
	}

	for _, test := range tests {
		t.Run("Test Jwt", func(t *testing.T) {
			s, err := utils.GenerateJwt(test.uR, configs.Cfg.AccessTokenExpireIn)
			assert.NoError(t, err)
			token, err := utils.VerifyJwt(s)
			assert.NoError(t, err)

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userId := int(claims["sub"].(float64))
				assert.Equal(t, test.uR.Id, userId)
			}

		})
	}
}
