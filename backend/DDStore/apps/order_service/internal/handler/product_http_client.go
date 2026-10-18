package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
)

type productHTTPClient struct {
	baseURL string
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

func (p productHTTPClient) UpdateProductState(productId string, state catalog.ProductState, ctx context.Context) (err error) {
	url := fmt.Sprintf("%s/", p.baseURL)

	jsonData, _ := json.Marshal(map[string]any{
		"product_id": productId,
		"state":      state,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update product state: unexpected status %d", resp.StatusCode)
	}

	return nil
}
