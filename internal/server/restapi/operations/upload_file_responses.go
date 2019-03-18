// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/onetwoclimb/internal/server/models"
)

// UploadFileOKCode is the HTTP code returned for type UploadFileOK
const UploadFileOKCode int = 200

/*UploadFileOK OK

swagger:response uploadFileOK
*/
type UploadFileOK struct {

	/*
	  In: Body
	*/
	Payload *UploadFileOKBody `json:"body,omitempty"`
}

// NewUploadFileOK creates UploadFileOK with default headers values
func NewUploadFileOK() *UploadFileOK {

	return &UploadFileOK{}
}

// WithPayload adds the payload to the upload file o k response
func (o *UploadFileOK) WithPayload(payload *UploadFileOKBody) *UploadFileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload file o k response
func (o *UploadFileOK) SetPayload(payload *UploadFileOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadFileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UploadFileInternalServerErrorCode is the HTTP code returned for type UploadFileInternalServerError
const UploadFileInternalServerErrorCode int = 500

/*UploadFileInternalServerError General server error

swagger:response uploadFileInternalServerError
*/
type UploadFileInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUploadFileInternalServerError creates UploadFileInternalServerError with default headers values
func NewUploadFileInternalServerError() *UploadFileInternalServerError {

	return &UploadFileInternalServerError{}
}

// WithPayload adds the payload to the upload file internal server error response
func (o *UploadFileInternalServerError) WithPayload(payload *models.Error) *UploadFileInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload file internal server error response
func (o *UploadFileInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadFileInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}