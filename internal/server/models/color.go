// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Color color object
// swagger:model Color
type Color struct {

	// color in hex
	// Required: true
	ColorHex *string `json:"colorHex"`

	// item name
	// Required: true
	Name *string `json:"name"`
}

// Validate validates this color
func (m *Color) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateColorHex(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Color) validateColorHex(formats strfmt.Registry) error {

	if err := validate.Required("colorHex", "body", m.ColorHex); err != nil {
		return err
	}

	return nil
}

func (m *Color) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Color) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Color) UnmarshalBinary(b []byte) error {
	var res Color
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
