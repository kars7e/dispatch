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
	swaggerclient "github.com/vmware/dispatch/pkg/secret-store/gen/client"
	secretclient "github.com/vmware/dispatch/pkg/secret-store/gen/client/secret"
	"github.com/vmware/dispatch/pkg/secret-store/gen/models"
)

// NO TESTS

// SecretsClient defines the secrets client interface
type SecretsClient interface {
	CreateSecret(context.Context, *Secret) (*Secret, error)
	DeleteSecret(ctx context.Context, secretName string) error
	UpdateSecret(context.Context, *Secret) (*Secret, error)
	GetSecret(ctx context.Context, secretName string) (*Secret, error)
	ListSecrets(context.Context) ([]Secret, error)
}

// Secret defines a secret
type Secret struct {
	models.Secret
}

// NewSecretsClient is used to create a new secrets client
func NewSecretsClient(host string, auth runtime.ClientAuthInfoWriter) *DefaultSecretsClient {
	transport := DefaultHTTPClient(host, swaggerclient.DefaultBasePath)
	return &DefaultSecretsClient{
		client: swaggerclient.New(transport, strfmt.Default),
		auth:   auth,
	}
}

// DefaultSecretsClient defines the default secrets client
type DefaultSecretsClient struct {
	client *swaggerclient.SecretStore
	auth   runtime.ClientAuthInfoWriter
}

// CreateSecret creates a secret
func (c *DefaultSecretsClient) CreateSecret(ctx context.Context, secret *Secret) (*Secret, error) {
	params := secretclient.AddSecretParams{
		Context: ctx,
		Secret:  &secret.Secret,
	}
	response, err := c.client.Secret.AddSecret(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when creating a secret")
	}
	return &Secret{Secret: *response.Payload}, nil
}

// DeleteSecret deletes a secret
func (c *DefaultSecretsClient) DeleteSecret(ctx context.Context, secretName string) error {
	params := secretclient.DeleteSecretParams{
		Context:    ctx,
		SecretName: secretName,
	}
	_, err := c.client.Secret.DeleteSecret(&params, c.auth)
	if err != nil {
		return errors.Wrap(err, "error when deleting a secret")
	}
	return nil
}

// UpdateSecret updates a secret
func (c *DefaultSecretsClient) UpdateSecret(ctx context.Context, secret *Secret) (*Secret, error) {
	params := secretclient.UpdateSecretParams{
		Context: ctx,
		Secret:  &secret.Secret,
	}
	response, err := c.client.Secret.UpdateSecret(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when updating a secret")
	}
	return &Secret{Secret: *response.Payload}, nil
}

// GetSecret retrieves a secret
func (c *DefaultSecretsClient) GetSecret(ctx context.Context, secretName string) (*Secret, error) {
	params := secretclient.GetSecretParams{
		Context:    ctx,
		SecretName: secretName,
	}
	response, err := c.client.Secret.GetSecret(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when retrieving a secret")
	}
	return &Secret{Secret: *response.Payload}, nil
}

// ListSecrets lists secrets
func (c *DefaultSecretsClient) ListSecrets(ctx context.Context) ([]Secret, error) {
	params := secretclient.GetSecretsParams{
		Context: ctx,
	}
	response, err := c.client.Secret.GetSecrets(&params, c.auth)
	if err != nil {
		return nil, errors.Wrap(err, "error when retrieving a secret")
	}
	var secrets []Secret
	for _, secret := range response.Payload {
		secrets = append(secrets, Secret{Secret: *secret})
	}
	return secrets, nil
}
