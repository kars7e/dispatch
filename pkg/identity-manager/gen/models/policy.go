///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

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

// Policy policy
// swagger:model Policy

type Policy struct {

	// created time
	// Read Only: true
	CreatedTime int64 `json:"createdTime,omitempty"`

	// id
	ID strfmt.UUID `json:"id,omitempty"`

	// modified time
	// Read Only: true
	ModifiedTime int64 `json:"modifiedTime,omitempty"`

	// name
	// Required: true
	// Pattern: ^[\w\d\-]+$
	Name *string `json:"name"`

	// rules
	// Required: true
	Rules PolicyRules `json:"rules"`

	// status
	// Read Only: true
	Status Status `json:"status,omitempty"`
}

/* polymorph Policy createdTime false */

/* polymorph Policy id false */

/* polymorph Policy modifiedTime false */

/* polymorph Policy name false */

/* polymorph Policy rules false */

/* polymorph Policy status false */

// Validate validates this policy
func (m *Policy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateRules(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Policy) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.Pattern("name", "body", string(*m.Name), `^[\w\d\-]+$`); err != nil {
		return err
	}

	return nil
}

func (m *Policy) validateRules(formats strfmt.Registry) error {

	if err := validate.Required("rules", "body", m.Rules); err != nil {
		return err
	}

	return nil
}

func (m *Policy) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Policy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Policy) UnmarshalBinary(b []byte) error {
	var res Policy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}