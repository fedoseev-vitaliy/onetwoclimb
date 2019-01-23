package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onetwoclimb/internal/server/restapi/operations"
)

func TestGetColors(t *testing.T) {
	h := New()
	req := operations.BoardColorsParams{
		HTTPRequest: &http.Request{},
	}

	r := h.GetColorsHandler(req)
	assert.NotNil(t, r)
}
