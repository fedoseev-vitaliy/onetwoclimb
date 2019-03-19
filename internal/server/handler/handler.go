package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/onetwoclimb/cmd/config"
	"github.com/onetwoclimb/internal/server/models"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/onetwoclimb/internal/storages"
	"github.com/sirupsen/logrus"
)

var l = logrus.New()

type Handler struct {
	MySQL  *storages.MySQLStorage
	config config.Config
}

func New(storage *storages.MySQLStorage, config config.Config) *Handler {
	return &Handler{MySQL: storage, config: config}
}

func (h *Handler) GetColorsHandler(params operations.GetBoardColorsParams) middleware.Responder {
	colors, err := h.MySQL.GetColors()
	if err != nil {
		l.WithError(err).Error("failed to get colors")
		return operations.NewGetBoardColorsInternalServerError()
	}

	res := make([]*models.Color, 0, len(colors))
	for _, color := range colors {
		res = append(res, &models.Color{
			ID:      color.Id,
			Name:    color.Name,
			PinCode: color.PinCode,
			Hex:     color.Hex,
		})
	}

	return operations.NewGetBoardColorsOK().WithPayload(&operations.GetBoardColorsOKBody{
		Colors: res,
	})
}

func (h *Handler) DeleteColorHandler(params operations.DelBoardColorParams) middleware.Responder {
	if err := h.MySQL.DelColor(int(params.ColorID)); err != nil {
		l.WithError(err).Error("failed to delete color")
		return operations.NewDelBoardColorInternalServerError()
	}
	return operations.NewDelBoardColorOK()
}

func (h *Handler) PostColorHandler(params operations.PostBoardColorsParams) middleware.Responder {
	if err := h.MySQL.PutColor(&storages.Color{
		Name:    params.Body.Name,
		PinCode: params.Body.PinCode,
		Hex:     params.Body.Hex,
	}); err != nil {
		l.WithError(err).Error("failed to add color")
		return operations.NewPostBoardColorsInternalServerError()
	}
	return operations.NewPostBoardColorsOK()
}

// todo refactor and move common logic for decod images to helpers
func ImageProducer() runtime.Producer {
	return runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		b, ok := data.(strfmt.Base64)
		if !ok {
			return errors.New("producer cast err")
		}

		ct := http.DetectContentType(b)
		switch {
		case strings.Contains(ct, "image/png"):
			img, err := png.Decode(bytes.NewReader(b))
			if err != nil {
				l.WithError(errors.WithStack(err)).Error("failed to decode png file in producer")
				return errors.WithStack(err)
			}

			// Prepare parent image where we want to position child image.
			target := image.NewRGBA(img.Bounds())
			// Draw child image.
			draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

			return png.Encode(writer, target)
		case strings.Contains(ct, "image/jpeg"):
			img, err := jpeg.Decode(bytes.NewReader(b))
			if err != nil {
				l.WithError(errors.WithStack(err)).Error("failed to decode png file in producer")
				return errors.WithStack(err)
			}

			// Prepare parent image where we want to position child image.
			target := image.NewRGBA(img.Bounds())
			// Draw child image.
			draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

			return jpeg.Encode(writer, target, nil)
		default:
			return errors.New(fmt.Sprintf("Unsupported content type: %s", ct))
		}
	})
}

func (h *Handler) ConfigureHandlers(api *operations.OneTwoClimbAPI) {
	api.Logger = l.Printf
	api.GetBoardColorsHandler = operations.GetBoardColorsHandlerFunc(h.GetColorsHandler)
	api.PostBoardColorsHandler = operations.PostBoardColorsHandlerFunc(h.PostColorHandler)
	api.DelBoardColorHandler = operations.DelBoardColorHandlerFunc(h.DeleteColorHandler)
	api.UploadFileHandler = operations.UploadFileHandlerFunc(h.PostUploadFile)
	api.DownloadFileHandler = operations.DownloadFileHandlerFunc(h.GetDownloadFile)
	api.ImagePngImageJpegProducer = ImageProducer()
}
