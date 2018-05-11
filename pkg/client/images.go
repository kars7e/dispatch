///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

package client

import (
	"context"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	swaggerclient "github.com/vmware/dispatch/pkg/image-manager/gen/client"
	baseimageclient "github.com/vmware/dispatch/pkg/image-manager/gen/client/base_image"
	imageclient "github.com/vmware/dispatch/pkg/image-manager/gen/client/image"
	"github.com/vmware/dispatch/pkg/image-manager/gen/models"
)

// NO TESTS

// ImagesClient defines the image client interface
type ImagesClient interface {
	// Images
	CreateImage(context.Context, *Image) (*Image, error)
	DeleteImage(ctx context.Context, imageName string) (*Image, error)
	UpdateImage(context.Context, *Image) (*Image, error)
	GetImage(ctx context.Context, imageName string) (*Image, error)
	ListImages(context.Context) ([]Image, error)

	// BaseImages
	CreateBaseImage(context.Context, *BaseImage) (*BaseImage, error)
	DeleteBaseImage(ctx context.Context, baseImageName string) (*BaseImage, error)
	UpdateBaseImage(context.Context, *BaseImage) (*BaseImage, error)
	GetBaseImage(ctx context.Context, baseImageName string) (*BaseImage, error)
	ListBaseImages(context.Context) ([]BaseImage, error)
}

// Image defines an image
type Image struct {
	models.Image
}

// BaseImage defines a base image
type BaseImage struct {
	models.BaseImage
}

// NewImagesClient is used to create a new Images client
func NewImagesClient(host string, auth runtime.ClientAuthInfoWriter) *DefaultImagesClient {
	transport := DefaultHTTPClient(host, swaggerclient.DefaultBasePath)
	return &DefaultImagesClient{
		client: swaggerclient.New(transport, strfmt.Default),
		auth:   auth,
	}
}

// DefaultImagesClient defines the default images client
type DefaultImagesClient struct {
	client *swaggerclient.ImageManager
	auth   runtime.ClientAuthInfoWriter
}

// CreateImage creates new image
func (c *DefaultImagesClient) CreateImage(ctx context.Context, image *Image) (*Image, error) {
	params := imageclient.AddImageParams{
		Context: ctx,
		Body:    &image.Image,
	}
	response, err := c.client.Image.AddImage(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when creating the image")
	}
	return &Image{Image: *response.Payload}, nil
}

// DeleteImage deletes an image
func (c *DefaultImagesClient) DeleteImage(ctx context.Context, imageName string) (*Image, error) {
	params := imageclient.DeleteImageByNameParams{
		Context:   ctx,
		ImageName: imageName,
	}
	response, err := c.client.Image.DeleteImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when deleting the image")
	}
	return &Image{Image: *response.Payload}, nil
}

// UpdateImage updates an image
func (c *DefaultImagesClient) UpdateImage(ctx context.Context, image *Image) (*Image, error) {
	params := imageclient.UpdateImageByNameParams{
		Context:   ctx,
		Body:      &image.Image,
		ImageName: *image.Name,
	}
	response, err := c.client.Image.UpdateImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when updating the image")
	}
	return &Image{Image: *response.Payload}, nil
}

// GetImage retrieves an image
func (c *DefaultImagesClient) GetImage(ctx context.Context, imageName string) (*Image, error) {
	params := imageclient.GetImageByNameParams{
		Context:   ctx,
		ImageName: imageName,
	}
	response, err := c.client.Image.GetImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the image")
	}
	return &Image{Image: *response.Payload}, nil
}

// ListImages returns a list of images
func (c *DefaultImagesClient) ListImages(ctx context.Context) ([]Image, error) {
	params := imageclient.GetImagesParams{
		Context: ctx,
	}
	response, err := c.client.Image.GetImages(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when listing images")
	}
	var images []Image
	for _, image := range response.Payload {
		images = append(images, Image{Image: *image})
	}
	return images, nil
}

// CreateBaseImage creates new base image
func (c *DefaultImagesClient) CreateBaseImage(ctx context.Context, image *BaseImage) (*BaseImage, error) {
	params := baseimageclient.AddBaseImageParams{
		Context: ctx,
		Body:    &image.BaseImage,
	}
	response, err := c.client.BaseImage.AddBaseImage(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when creating the base image")
	}
	return &BaseImage{BaseImage: *response.Payload}, nil
}

// DeleteBaseImage deletes the base image
func (c *DefaultImagesClient) DeleteBaseImage(ctx context.Context, baseImageName string) (*BaseImage, error) {
	params := baseimageclient.DeleteBaseImageByNameParams{
		Context:       ctx,
		BaseImageName: baseImageName,
	}
	response, err := c.client.BaseImage.DeleteBaseImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when deleting the base image")
	}
	return &BaseImage{BaseImage: *response.Payload}, nil
}

// UpdateBaseImage updates the base image
func (c *DefaultImagesClient) UpdateBaseImage(ctx context.Context, image *BaseImage) (*BaseImage, error) {
	params := baseimageclient.UpdateBaseImageByNameParams{
		Context:       ctx,
		Body:          &image.BaseImage,
		BaseImageName: *image.Name,
	}
	response, err := c.client.BaseImage.UpdateBaseImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when updating the base image")
	}
	return &BaseImage{BaseImage: *response.Payload}, nil
}

// GetBaseImage retrieves the base image
func (c *DefaultImagesClient) GetBaseImage(ctx context.Context, baseImageName string) (*BaseImage, error) {
	params := baseimageclient.GetBaseImageByNameParams{
		Context:       ctx,
		BaseImageName: baseImageName,
	}
	response, err := c.client.BaseImage.GetBaseImageByName(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when retreiving the base image")
	}
	return &BaseImage{BaseImage: *response.Payload}, nil
}

// ListBaseImages returns a list of base images
func (c *DefaultImagesClient) ListBaseImages(ctx context.Context) ([]BaseImage, error) {
	params := baseimageclient.GetBaseImagesParams{
		Context: ctx,
	}
	response, err := c.client.BaseImage.GetBaseImages(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when listing base images")
	}
	var images []BaseImage
	for _, image := range response.Payload {
		images = append(images, BaseImage{BaseImage: *image})
	}
	return images, nil
}
