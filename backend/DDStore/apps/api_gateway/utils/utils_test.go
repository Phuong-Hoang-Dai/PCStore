package utils_test

import (
	"testing"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/utils"
)

func TestParseCustomDuration(t *testing.T) {
	tests := []struct {
		s       string
		result  time.Duration
		wantErr bool
	}{
		{"20d", time.Duration(20*24) * time.Hour, false},
		{"200h", time.Duration(200) * time.Hour, false},
		{"2z", 0, true},
		{"2hd", 0, true},
		{"sdgsdgds", 0, true},
	}

	for _, test := range tests {
		t.Run("Test ParseCustomDuration", func(t *testing.T) {
			result, err := utils.ParseCustomDuration(test.s)
			if (err != nil) != test.wantErr {
				t.Errorf("Input: %v \n Error : %v, Result: %v", test, err, result)
			}
			if result != test.result {
				t.Errorf("Input: %v \n Error : %v, Result: %v", test, err, result)
			}
		})
	}
}

func TestSpaceToDash(t *testing.T) {
	tests := []struct {
		s       string
		result  string
		wantErr bool
	}{
		{s: "PC Store", result: "PC-Store", wantErr: false},
	}

	for _, test := range tests {
		t.Run("Test SpaceToDash", func(t *testing.T) {
			result := utils.SpaceToDash(test.s)
			if result != test.result {
				t.Errorf("Input: %v \n Result: %v", test, result)
			}
		})
	}
}
