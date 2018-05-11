///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////
package injectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vmware/dispatch/pkg/client"
	"github.com/vmware/dispatch/pkg/client/mocks"
	"github.com/vmware/dispatch/pkg/functions"
	"github.com/vmware/dispatch/pkg/secret-store/gen/models"
)

//go:generate mockery -name SecretInjector -case underscore -dir .

func TestInjectSecret(t *testing.T) {

	expectedSecretName := "testSecret"
	expectedSecretValue := models.SecretValue{"secret1": "value1", "secret2": "value2"}

	expectedOutput := map[string]interface{}{"secret1": "value1", "secret2": "value2"}

	secret := client.Secret{}
	secretsClient := &mocks.SecretsClient{}
	secretsClient.On("GetSecret", mock.Anything, mock.Anything).Return(
		&client.Secret{
			Secret: models.Secret{
				Name:    &expectedSecretName,
				Secrets: expectedSecretValue,
			}}, nil)

	injector := NewSecretInjector(secretsClient)

	cookie := "testCookie"

	printSecretsFn := func(ctx functions.Context, _ interface{}) (interface{}, error) {
		return ctx["secrets"], nil
	}

	ctx := functions.Context{}
	output, err := injector.GetMiddleware([]string{expectedSecretName}, cookie)(printSecretsFn)(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}
