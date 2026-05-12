package utils

import (
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(user model.UseResp, expConfig string) (string, error) {
	key := []byte(configs.Cfg.JWTSecret)

	iss := SpaceToDash(configs.Cfg.AppName)
	exp, err := ParseCustomDuration(expConfig)
	if err != nil {
		return "", err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Id,
		"iss":  "Server-" + iss,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(exp).Unix(),
		"role": user.Role,
	})

	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}

func VerifyJwt(s string) (*jwt.Token, error) {
	key := []byte(configs.Cfg.JWTSecret)

	token, err := jwt.Parse(s, func(t *jwt.Token) (any, error) {
		return key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RefreshToken(s string) (string, error) {
	token, err := VerifyJwt(s)
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId := int(claims["sub"].(float64))
		u := model.UseResp{Id: userId, Role: claims["role"].(string)}

		if r, err := GenerateJwt(u, configs.Cfg.AccessTokenExpireIn); err != nil {
			return "", err
		} else {
			return r, nil
		}
	} else {
		return "", Err.ErrCannotReadRefreshToken
	}
}
