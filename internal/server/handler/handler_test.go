package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onetwoclimb/internal/server/restapi/operations"
)

func TestGetColorsCount(t *testing.T) {
	h := New()
	req := operations.BoardColorsParams{
		HTTPRequest: &http.Request{},
	}

	r := h.GetColorsHandler(req)
	res := r.(*operations.BoardColorsOK)
	assert.NotNil(t, res)
	assert.Equal(t, 7, len(res.Payload.Colors))
}
