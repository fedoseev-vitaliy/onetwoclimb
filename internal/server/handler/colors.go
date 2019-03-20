package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/onetwoclimb/internal/server/models"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/onetwoclimb/internal/storages"
)

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
