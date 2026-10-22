package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
)

type userHTTPClient struct {
	baseURL string
}

func (u userHTTPClient) GetUserById(userId int, ctx context.Context) (user service.UserDTO, err error) {
	url := fmt.Sprintf("%s/%d", u.baseURL, userId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return user, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return user, errors.New("failed to read body respone")
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg string
		if rawErr, ok := raw["error"]; ok {
			json.Unmarshal(rawErr, &errMsg)
		}
		return user, fmt.Errorf("get user info: unexpected status %d: %s", resp.StatusCode, errMsg)
	}

	rawData, ok := raw["data"]
	if !ok {
		return user, errors.New("failed to get data")
	}
	if err := json.Unmarshal(rawData, &user); err != nil {
		return user, err
	}

	return user, nil
}
