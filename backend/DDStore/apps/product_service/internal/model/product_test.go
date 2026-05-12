package model_test

import (
	"testing"

	model "github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPagingProcess(t *testing.T) {
	tests := []struct {
		p              model.Paging
		expectedLimit  int
		expectedOffset int
	}{
		{model.Paging{Limit: 0, Offset: 1}, 0, 1},
		{model.Paging{Limit: -1, Offset: 1}, 0, 1},
		{model.Paging{Limit: 0, Offset: -1}, 0, 0},
		{model.Paging{Limit: 100, Offset: 0}, 50, 0},
		{model.Paging{Limit: -1, Offset: -1}, 0, 0},
	}

	for _, test := range tests {
		t.Run("TestCase", func(t *testing.T) {

			test.p.Process()
			assert.Equal(t, test.expectedLimit, test.p.Limit)
			assert.Equal(t, test.expectedOffset, test.p.Offset)
		})
	}
}
