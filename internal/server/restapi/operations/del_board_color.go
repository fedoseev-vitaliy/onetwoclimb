// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DelBoardColorHandlerFunc turns a function with the right signature into a del board color handler
type DelBoardColorHandlerFunc func(DelBoardColorParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DelBoardColorHandlerFunc) Handle(params DelBoardColorParams) middleware.Responder {
	return fn(params)
}

// DelBoardColorHandler interface for that can handle valid del board color params
type DelBoardColorHandler interface {
	Handle(DelBoardColorParams) middleware.Responder
}

// NewDelBoardColor creates a new http.Handler for the del board color operation
func NewDelBoardColor(ctx *middleware.Context, handler DelBoardColorHandler) *DelBoardColor {
	return &DelBoardColor{Context: ctx, Handler: handler}
}

/*DelBoardColor swagger:route DELETE /colors/{colorId} delBoardColor

delete borad color

*/
type DelBoardColor struct {
	Context *middleware.Context
	Handler DelBoardColorHandler
}

func (o *DelBoardColor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDelBoardColorParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
