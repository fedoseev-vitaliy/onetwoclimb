package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onetwoclimb/cmd/config"
	"github.com/onetwoclimb/internal/server/restapi/operations"
)

func TestGetColorsCount(t *testing.T) {
	t.Skip("need to setup dockertest")
	h := New(nil, config.Config{MaxFileSize: 5000000})
	req := operations.GetBoardColorsParams{
		HTTPRequest: &http.Request{},
	}

	r := h.GetColorsHandler(req)
	res := r.(*operations.GetBoardColorsOK)
	assert.NotNil(t, res)
	assert.Equal(t, 7, len(res.Payload.Colors))
}
