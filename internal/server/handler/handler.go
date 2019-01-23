package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/onetwoclimb/internal/server/models"
	"github.com/onetwoclimb/internal/server/restapi/operations"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) GetColorsHandler(parameters operations.BoardColorsParams) middleware.Responder {
	return operations.NewBoardColorsOK().WithPayload(&operations.BoardColorsOKBody{
		Colors: h.getColors(),
	})
}

func (h *Handler) getColors() []*models.Color {
	return []*models.Color{
		{ID: 1, Name: "start", PinCode: "001", Hex: "#0040ff"},
		{ID: 2, Name: "finish", PinCode: "100", Hex: "#ff0000"},
		{ID: 3, Name: "route", PinCode: "010", Hex: "#04ff00"},
		{ID: 4, Name: "blank", PinCode: "000", Hex: ""},
		{ID: 5, Name: "event_flash", PinCode: "", Hex: "#FFA726"},
		{ID: 6, Name: "event_top", PinCode: "", Hex: "#008BA3"},
		{ID: 7, Name: "event_zone", PinCode: "", Hex: "#00BCD4"},
	}
}

func (h *Handler) ConfigureHandlers(api *operations.OneTwoClimbAPI) {
	api.BoardColorsHandler = operations.BoardColorsHandlerFunc(h.GetColorsHandler)
}
