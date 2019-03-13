package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/onetwoclimb/cmd/config"
	"github.com/onetwoclimb/internal/server/models"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/onetwoclimb/internal/storages"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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

func (h *Handler) PostUploadFile(params operations.UploadFileParams) middleware.Responder {
	defer params.File.Close()

	lr := &io.LimitedReader{N: int64(h.config.MaxFileSize), R: params.File}
	ib, err := ioutil.ReadAll(lr)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to read image")
		return operations.NewUploadFileInternalServerError()
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to generate filename")
		return operations.NewUploadFileInternalServerError()
	}
	fileName := uuid.String()
	// todo add other img types and think how to reuse code below
	ct := http.DetectContentType(ib)
	switch {
	case strings.Contains(ct, "image/png"):
		img, err := png.Decode(bytes.NewReader(ib))
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to decode png file")
			return operations.NewUploadFileInternalServerError()
		}

		// Prepare parent image where we want to position child image.
		target := image.NewRGBA(img.Bounds())
		// Draw child image.
		draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

		filePath := fmt.Sprintf("%s/%s.png", h.config.FilesDst, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to create png file")
			return operations.NewUploadFileInternalServerError()
		}

		err = png.Encode(f, target)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to save file")
			f.Close()
			if err := os.Remove(filePath); err != nil {
				l.WithError(errors.WithStack(err)).Errorf("failed to delete file by path: %s", filePath)
			}
			return operations.NewUploadFileInternalServerError()
		}
	default:
		l.WithError(errors.New("unsupported content type")).Errorf("content type:%s", ct)
		return operations.NewUploadFileInternalServerError()
	}

	return operations.NewUploadFileOK().WithPayload(&operations.UploadFileOKBody{ID: &fileName})
}

// todo implement download
//func (h *Handler) GetDownloadFile(params operations.DownloadFileParams) middleware.Responder {
//
//	return operations.NewDownloadFileOK().WithPayload()
//}

func (h *Handler) ConfigureHandlers(api *operations.OneTwoClimbAPI) {
	api.Logger = l.Printf
	api.GetBoardColorsHandler = operations.GetBoardColorsHandlerFunc(h.GetColorsHandler)
	api.PostBoardColorsHandler = operations.PostBoardColorsHandlerFunc(h.PostColorHandler)
	api.DelBoardColorHandler = operations.DelBoardColorHandlerFunc(h.DeleteColorHandler)
	api.UploadFileHandler = operations.UploadFileHandlerFunc(h.PostUploadFile)
}
