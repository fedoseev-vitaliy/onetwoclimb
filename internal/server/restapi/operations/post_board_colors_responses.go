// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/onetwoclimb/internal/server/models"
)

// PostBoardColorsOKCode is the HTTP code returned for type PostBoardColorsOK
const PostBoardColorsOKCode int = 200

/*PostBoardColorsOK OK

swagger:response postBoardColorsOK
*/
type PostBoardColorsOK struct {
}

// NewPostBoardColorsOK creates PostBoardColorsOK with default headers values
func NewPostBoardColorsOK() *PostBoardColorsOK {

	return &PostBoardColorsOK{}
}

// WriteResponse to the client
func (o *PostBoardColorsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PostBoardColorsInternalServerErrorCode is the HTTP code returned for type PostBoardColorsInternalServerError
const PostBoardColorsInternalServerErrorCode int = 500

/*PostBoardColorsInternalServerError General server error

swagger:response postBoardColorsInternalServerError
*/
type PostBoardColorsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostBoardColorsInternalServerError creates PostBoardColorsInternalServerError with default headers values
func NewPostBoardColorsInternalServerError() *PostBoardColorsInternalServerError {

	return &PostBoardColorsInternalServerError{}
}

// WithPayload adds the payload to the post board colors internal server error response
func (o *PostBoardColorsInternalServerError) WithPayload(payload *models.Error) *PostBoardColorsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post board colors internal server error response
func (o *PostBoardColorsInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBoardColorsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
