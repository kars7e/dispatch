///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

// Code generated by go-swagger; DO NOT EDIT.

package organization

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/vmware/dispatch/pkg/api/v1"
)

// UpdateOrganizationReader is a Reader for the UpdateOrganization structure.
type UpdateOrganizationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateOrganizationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUpdateOrganizationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewUpdateOrganizationBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewUpdateOrganizationNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewUpdateOrganizationInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateOrganizationOK creates a UpdateOrganizationOK with default headers values
func NewUpdateOrganizationOK() *UpdateOrganizationOK {
	return &UpdateOrganizationOK{}
}

/*UpdateOrganizationOK handles this case with default header values.

Successful update
*/
type UpdateOrganizationOK struct {
	Payload *v1.Organization
}

func (o *UpdateOrganizationOK) Error() string {
	return fmt.Sprintf("[PUT /v1/iam/organization/{organizationName}][%d] updateOrganizationOK  %+v", 200, o.Payload)
}

func (o *UpdateOrganizationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(v1.Organization)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateOrganizationBadRequest creates a UpdateOrganizationBadRequest with default headers values
func NewUpdateOrganizationBadRequest() *UpdateOrganizationBadRequest {
	return &UpdateOrganizationBadRequest{}
}

/*UpdateOrganizationBadRequest handles this case with default header values.

Invalid input
*/
type UpdateOrganizationBadRequest struct {
	Payload *v1.Error
}

func (o *UpdateOrganizationBadRequest) Error() string {
	return fmt.Sprintf("[PUT /v1/iam/organization/{organizationName}][%d] updateOrganizationBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateOrganizationBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(v1.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateOrganizationNotFound creates a UpdateOrganizationNotFound with default headers values
func NewUpdateOrganizationNotFound() *UpdateOrganizationNotFound {
	return &UpdateOrganizationNotFound{}
}

/*UpdateOrganizationNotFound handles this case with default header values.

Organization not found
*/
type UpdateOrganizationNotFound struct {
	Payload *v1.Error
}

func (o *UpdateOrganizationNotFound) Error() string {
	return fmt.Sprintf("[PUT /v1/iam/organization/{organizationName}][%d] updateOrganizationNotFound  %+v", 404, o.Payload)
}

func (o *UpdateOrganizationNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(v1.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateOrganizationInternalServerError creates a UpdateOrganizationInternalServerError with default headers values
func NewUpdateOrganizationInternalServerError() *UpdateOrganizationInternalServerError {
	return &UpdateOrganizationInternalServerError{}
}

/*UpdateOrganizationInternalServerError handles this case with default header values.

Internal error
*/
type UpdateOrganizationInternalServerError struct {
	Payload *v1.Error
}

func (o *UpdateOrganizationInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /v1/iam/organization/{organizationName}][%d] updateOrganizationInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateOrganizationInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(v1.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
