// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/onetwoclimb/internal/server/models"
)

// DelBoardColorOKCode is the HTTP code returned for type DelBoardColorOK
const DelBoardColorOKCode int = 200

/*DelBoardColorOK OK

swagger:response delBoardColorOK
*/
type DelBoardColorOK struct {
}

// NewDelBoardColorOK creates DelBoardColorOK with default headers values
func NewDelBoardColorOK() *DelBoardColorOK {

	return &DelBoardColorOK{}
}

// WriteResponse to the client
func (o *DelBoardColorOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DelBoardColorInternalServerErrorCode is the HTTP code returned for type DelBoardColorInternalServerError
const DelBoardColorInternalServerErrorCode int = 500

/*DelBoardColorInternalServerError General server error. Error codes:
  - 4 Server error

swagger:response delBoardColorInternalServerError
*/
type DelBoardColorInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDelBoardColorInternalServerError creates DelBoardColorInternalServerError with default headers values
func NewDelBoardColorInternalServerError() *DelBoardColorInternalServerError {

	return &DelBoardColorInternalServerError{}
}

// WithPayload adds the payload to the del board color internal server error response
func (o *DelBoardColorInternalServerError) WithPayload(payload *models.Error) *DelBoardColorInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the del board color internal server error response
func (o *DelBoardColorInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DelBoardColorInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
