///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////
package injectors

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vmware/dispatch/pkg/client"
	clientmocks "github.com/vmware/dispatch/pkg/client/mocks"
	"github.com/vmware/dispatch/pkg/functions"
	"github.com/vmware/dispatch/pkg/functions/mocks"
	"github.com/vmware/dispatch/pkg/secret-store/gen/models"
	serviceclient "github.com/vmware/dispatch/pkg/service-manager/gen/client"
	service "github.com/vmware/dispatch/pkg/service-manager/gen/client/service_instance"
	servicemodels "github.com/vmware/dispatch/pkg/service-manager/gen/models"
)

//go:generate mockery -name ServiceInjector -case underscore -dir .

func TestInjectService(t *testing.T) {

	expectedSecretValue := models.SecretValue{"secret1": "value1", "secret2": "value2"}
	expectedServiceName := "testService"
	expectedOutput := map[string]interface{}{"secret1": "value1", "secret2": "value2"}

	serviceID := uuid.NewV4().String()

	serviceTransport := &mocks.ClientTransport{}
	serviceTransport.On("Submit", mock.Anything).Return(
		&service.GetServiceInstanceByNameOK{
			Payload: &servicemodels.ServiceInstance{
				Name: &expectedServiceName,
				ID:   strfmt.UUID(serviceID),
				Binding: &servicemodels.ServiceBinding{
					Status: servicemodels.StatusREADY,
				},
			}}, nil)

	secretsClient := &clientmocks.SecretsClient{}
	secretsClient.On("GetSecret", mock.Anything, mock.Anything).Return(
		&client.Secret{
			Secret: models.Secret{
				Name:    &serviceID,
				Secrets: expectedSecretValue,
			}}, nil)

	serviceManager := serviceclient.New(serviceTransport, strfmt.Default)

	injector := NewServiceInjector(secretsClient, serviceManager)

	cookie := "testCookie"

	printServiceFn := func(ctx functions.Context, _ interface{}) (interface{}, error) {
		return ctx["serviceBindings"].(map[string]interface{})[expectedServiceName], nil
	}

	ctx := functions.Context{}
	output, err := injector.GetMiddleware([]string{expectedServiceName}, cookie)(printServiceFn)(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}
