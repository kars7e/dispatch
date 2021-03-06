///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetVersionHandlerFunc turns a function with the right signature into a get version handler
type GetVersionHandlerFunc func(GetVersionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetVersionHandlerFunc) Handle(params GetVersionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetVersionHandler interface for that can handle valid get version params
type GetVersionHandler interface {
	Handle(GetVersionParams, interface{}) middleware.Responder
}

// NewGetVersion creates a new http.Handler for the get version operation
func NewGetVersion(ctx *middleware.Context, handler GetVersionHandler) *GetVersion {
	return &GetVersion{Context: ctx, Handler: handler}
}

/*GetVersion swagger:route GET /v1/version getVersion

get version info

*/
type GetVersion struct {
	Context *middleware.Context
	Handler GetVersionHandler
}

func (o *GetVersion) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetVersionParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
