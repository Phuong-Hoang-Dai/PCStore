package utils

import (
	"strconv"
	"strings"
	"time"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
)

func ParseCustomDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, Err.ErrTryToParseEmptyStringToTime
	}
	duration, err := strconv.Atoi(s[0 : len(s)-1])
	if err != nil {
		return 0, err
	}
	var durationTime time.Duration
	switch suffix := s[len(s)-1]; suffix {
	case 'h':
		durationTime = time.Duration(duration) * time.Hour
	case 'd':
		durationTime = time.Duration(duration) * time.Hour * 24
	default:
		return 0, Err.ErrTrytoParseWrongFormatStringToTime
	}
	return durationTime, nil
}

func SpaceToDash(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
