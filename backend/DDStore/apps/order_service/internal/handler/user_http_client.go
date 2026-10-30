package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
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

	_, err = callUserService(req, &user)
	return user, err
}

// GetOrCreateGuest posts the guest without a password, which the user service treats as a guest:
// it returns the account using the email, or creates one that can't log in.
func (u userHTTPClient) GetOrCreateGuest(guest service.GuestDTO, ctx context.Context) (int, error) {
	body, err := json.Marshal(guest)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.baseURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	var data struct {
		Id int `json:"id"`
	}
	status, err := callUserService(req, &data)
	if status == http.StatusBadRequest {
		return 0, fmt.Errorf("%w: %v", model.ErrInvalidGuest, err)
	}
	return data.Id, err
}

// callUserService sends req and decodes the "data" field of a successful response into out.
func callUserService(req *http.Request, out any) (status int, err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resp.StatusCode, errors.New("failed to read body respone")
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg string
		if rawErr, ok := raw["error"]; ok {
			json.Unmarshal(rawErr, &errMsg)
		}
		return resp.StatusCode, fmt.Errorf("user service: unexpected status %d: %s", resp.StatusCode, errMsg)
	}

	rawData, ok := raw["data"]
	if !ok {
		return resp.StatusCode, errors.New("failed to get data")
	}
	return resp.StatusCode, json.Unmarshal(rawData, out)
}
