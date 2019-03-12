package handler

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/onetwoclimb/internal/storages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/onetwoclimb/internal/server/models"
	"github.com/onetwoclimb/internal/server/restapi/operations"
)

var l = logrus.New()

type Handler struct {
	MySQL *storages.MySQLStorage
}

func New(storage *storages.MySQLStorage) *Handler {
	return &Handler{MySQL: storage}
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

func (h *Handler) PostUpload(params operations.UploadFileParams) middleware.Responder {
	img, err := png.Decode(params.File)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to decode png file")
		return operations.NewUploadFileInternalServerError()
	}

	f, err := os.Create("img.png") // todo make unique name and dest path to config
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to create png file")
		return operations.NewUploadFileInternalServerError()
	}
	defer f.Close()

	// Prepare parent image where we want to position child image.
	target := image.NewRGBA(img.Bounds())
	// Draw child image.
	draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

	err = png.Encode(f, target)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to save file")
		return operations.NewUploadFileInternalServerError()
	}

	return operations.NewUploadFileOK()
}

func (h *Handler) ConfigureHandlers(api *operations.OneTwoClimbAPI) {
	api.Logger = l.Printf
	api.GetBoardColorsHandler = operations.GetBoardColorsHandlerFunc(h.GetColorsHandler)
	api.PostBoardColorsHandler = operations.PostBoardColorsHandlerFunc(h.PostColorHandler)
	api.DelBoardColorHandler = operations.DelBoardColorHandlerFunc(h.DeleteColorHandler)
	api.UploadFileHandler = operations.UploadFileHandlerFunc(h.PostUpload)
}
