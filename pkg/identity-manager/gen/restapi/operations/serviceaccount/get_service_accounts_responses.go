///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

// Code generated by go-swagger; DO NOT EDIT.

package serviceaccount

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/vmware/dispatch/pkg/api/v1"
)

// GetServiceAccountsOKCode is the HTTP code returned for type GetServiceAccountsOK
const GetServiceAccountsOKCode int = 200

/*GetServiceAccountsOK Successful operation

swagger:response getServiceAccountsOK
*/
type GetServiceAccountsOK struct {

	/*
	  In: Body
	*/
	Payload []*v1.ServiceAccount `json:"body,omitempty"`
}

// NewGetServiceAccountsOK creates GetServiceAccountsOK with default headers values
func NewGetServiceAccountsOK() *GetServiceAccountsOK {

	return &GetServiceAccountsOK{}
}

// WithPayload adds the payload to the get service accounts o k response
func (o *GetServiceAccountsOK) WithPayload(payload []*v1.ServiceAccount) *GetServiceAccountsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service accounts o k response
func (o *GetServiceAccountsOK) SetPayload(payload []*v1.ServiceAccount) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceAccountsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*v1.ServiceAccount, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetServiceAccountsInternalServerErrorCode is the HTTP code returned for type GetServiceAccountsInternalServerError
const GetServiceAccountsInternalServerErrorCode int = 500

/*GetServiceAccountsInternalServerError Internal Error

swagger:response getServiceAccountsInternalServerError
*/
type GetServiceAccountsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *v1.Error `json:"body,omitempty"`
}

// NewGetServiceAccountsInternalServerError creates GetServiceAccountsInternalServerError with default headers values
func NewGetServiceAccountsInternalServerError() *GetServiceAccountsInternalServerError {

	return &GetServiceAccountsInternalServerError{}
}

// WithPayload adds the payload to the get service accounts internal server error response
func (o *GetServiceAccountsInternalServerError) WithPayload(payload *v1.Error) *GetServiceAccountsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service accounts internal server error response
func (o *GetServiceAccountsInternalServerError) SetPayload(payload *v1.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceAccountsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetServiceAccountsDefault Unexpected Error

swagger:response getServiceAccountsDefault
*/
type GetServiceAccountsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *v1.Error `json:"body,omitempty"`
}

// NewGetServiceAccountsDefault creates GetServiceAccountsDefault with default headers values
func NewGetServiceAccountsDefault(code int) *GetServiceAccountsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetServiceAccountsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get service accounts default response
func (o *GetServiceAccountsDefault) WithStatusCode(code int) *GetServiceAccountsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get service accounts default response
func (o *GetServiceAccountsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get service accounts default response
func (o *GetServiceAccountsDefault) WithPayload(payload *v1.Error) *GetServiceAccountsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service accounts default response
func (o *GetServiceAccountsDefault) SetPayload(payload *v1.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceAccountsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
