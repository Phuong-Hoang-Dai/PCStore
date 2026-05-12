package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

func makeRequest(method string, url string, body io.Reader) (map[string]json.RawMessage, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, Err.ErrCannotReadBody
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, Err.ErrCannotReadBody
	}
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, resp.StatusCode, errors.New("failed to read body")
	}

	return raw, resp.StatusCode, nil
}

func createData(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
		}
		bodyData, statusCode, err := makeRequest(http.MethodPost, baseUrl, bytes.NewBuffer(body))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func getDatas(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		bodyData, statusCode, err := makeRequest(http.MethodGet, fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, offset), bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}
func getDataById(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		bodyData, statusCode, err := makeRequest(http.MethodGet, fmt.Sprintf("%s/%d", baseUrl, id), bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func updateDataById(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
		}
		bodyData, statusCode, err := makeRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseUrl, id), bytes.NewBuffer(body))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func deleteDataById(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		bodyData, statusCode, err := makeRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseUrl, id), bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
