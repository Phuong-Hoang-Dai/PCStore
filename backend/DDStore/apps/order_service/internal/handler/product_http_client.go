package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
)

type productHTTPClient struct {
	baseURL string
}

func (p productHTTPClient) GetStock(items []service.OrderDTO) error {
	url := fmt.Sprintf("%s/getstock", p.baseURL)

	jsonData, _ := json.Marshal(items)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(body, &raw); err != nil {
			return errors.New("failed to read body respone")
		}

		var errMsg string
		if rawErr, ok := raw["error"]; ok {
			json.Unmarshal(rawErr, &errMsg)
			return errors.New(errMsg)
		}

		return errors.New("failed to get stock")
	}

	return nil
}

func (p productHTTPClient) RestoreStock(items []service.OrderDTO) error {
	url := fmt.Sprintf("%s/restore", p.baseURL)

	jsonData, _ := json.Marshal(items)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(body, &raw); err != nil {
			return errors.New("failed to read body respone")
		}

		var errMsg string
		if rawErr, ok := raw["error"]; ok {
			json.Unmarshal(rawErr, &errMsg)
			return errors.New(errMsg)
		}

		return errors.New("failed to get stock")
	}

	return nil
}

func (p productHTTPClient) GetPriceProduct(items *[]service.OrderDTO) error {
	url := fmt.Sprintf("%s/getprice", p.baseURL)

	jsonData, _ := json.Marshal(items)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return errors.New("failed to read body respone")
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg string
		if rawErr, ok := raw["error"]; ok {
			json.Unmarshal(rawErr, &errMsg)
		}
		return errors.New("failed to get price product")
	}

	if rawData, ok := raw["data"]; ok {
		var oR []service.OrderResponeDTO
		json.Unmarshal(rawData, &oR)
		if len(oR) != len(*items) {
			return errors.New("something wrong happend")
		}
		for i := range *items {
			service.MapOrderResponeDTOtoOrderDTO(oR[i], &(*items)[i])
		}
	} else {
		return errors.New("failed to get data")
	}

	return nil
}
