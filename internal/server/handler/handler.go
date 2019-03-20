package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/onetwoclimb/cmd/config"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/onetwoclimb/internal/storages"
)

// todo improve logger and panic handler with pretty output
// todo move everything to docker

var l = logrus.New()

type Handler struct {
	MySQL  *storages.MySQLStorage
	config config.Config
}

func New(storage *storages.MySQLStorage, config config.Config) *Handler {
	return &Handler{MySQL: storage, config: config}
}

func ImageProducer() runtime.Producer {
	return runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		b, ok := data.(strfmt.Base64)
		if !ok {
			return errors.New("producer cast err")
		}

		ct := http.DetectContentType(b)
		switch {
		case strings.Contains(ct, "image/png"):
			return EncodePNG(b, writer)
		case strings.Contains(ct, "image/jpeg"):
			return EncodeJPEG(b, writer)
		default:
			return errors.New(fmt.Sprintf("Unsupported content type: %s", ct))
		}
	})
}

// todo add internal error codes
func (h *Handler) ConfigureHandlers(api *operations.OneTwoClimbAPI) {
	api.Logger = l.Printf
	api.GetBoardColorsHandler = operations.GetBoardColorsHandlerFunc(h.GetColorsHandler)
	api.PostBoardColorsHandler = operations.PostBoardColorsHandlerFunc(h.PostColorHandler)
	api.DelBoardColorHandler = operations.DelBoardColorHandlerFunc(h.DeleteColorHandler)
	api.UploadFileHandler = operations.UploadFileHandlerFunc(h.PostUploadFile)
	api.DownloadFileHandler = operations.DownloadFileHandlerFunc(h.GetDownloadFile)
	api.ImagePngImageJpegProducer = ImageProducer()
}
