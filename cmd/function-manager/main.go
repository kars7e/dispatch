///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

package main

import (
	"os"
	"time"

	docker "github.com/docker/docker/client"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/swag"
	"github.com/jessevdk/go-flags"
	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/dispatch/pkg/client"

	"github.com/vmware/dispatch/pkg/config"
	"github.com/vmware/dispatch/pkg/entity-store"
	"github.com/vmware/dispatch/pkg/function-manager"
	"github.com/vmware/dispatch/pkg/function-manager/gen/restapi"
	"github.com/vmware/dispatch/pkg/function-manager/gen/restapi/operations"
	"github.com/vmware/dispatch/pkg/functions"
	"github.com/vmware/dispatch/pkg/functions/injectors"
	"github.com/vmware/dispatch/pkg/functions/kubeless"
	"github.com/vmware/dispatch/pkg/functions/noop"
	"github.com/vmware/dispatch/pkg/functions/openfaas"
	"github.com/vmware/dispatch/pkg/functions/openwhisk"
	"github.com/vmware/dispatch/pkg/functions/riff"
	"github.com/vmware/dispatch/pkg/functions/runner"
	"github.com/vmware/dispatch/pkg/functions/validator"
	"github.com/vmware/dispatch/pkg/middleware"
	"github.com/vmware/dispatch/pkg/utils"
)

var drivers = map[string]func() functions.FaaSDriver{
	"openfaas": func() functions.FaaSDriver {
		faas, err := openfaas.New(&openfaas.Config{
			Gateway:             config.Global.Function.OpenFaas.Gateway,
			K8sConfig:           config.Global.Function.OpenFaas.K8sConfig,
			FuncNamespace:       config.Global.Function.OpenFaas.FuncNamespace,
			FuncDefaultRequests: config.Global.Function.OpenFaas.FuncDefaultRequests,
			FuncDefaultLimits:   config.Global.Function.OpenFaas.FuncDefaultLimits,
			ImagePullSecret:     config.Global.Function.OpenFaas.ImagePullSecret,
		})
		if err != nil {
			log.Fatalf("Error starting OpenFaaS driver: %+v", err)
		}
		return faas
	},
	"riff": func() functions.FaaSDriver {
		faas, err := riff.New(&riff.Config{
			KafkaBrokers:        config.Global.Function.Riff.KafkaBrokers,
			K8sConfig:           config.Global.Function.Riff.K8sConfig,
			FuncNamespace:       config.Global.Function.Riff.FuncNamespace,
			FuncDefaultRequests: config.Global.Function.Riff.FuncDefaultRequests,
			FuncDefaultLimits:   config.Global.Function.Riff.FuncDefaultLimits,
		})
		if err != nil {
			log.Fatalf("Error starting riff driver: %+v", err)
		}
		return faas
	},
	"openwhisk": func() functions.FaaSDriver {
		faas, err := openwhisk.New(&openwhisk.Config{
			AuthToken: config.Global.Function.Openwhisk.AuthToken,
			Host:      config.Global.Function.Openwhisk.Host,
			Insecure:  true,
		})
		if err != nil {
			log.Fatalf("Error getting OpenWhisk client: %+v", err)
		}
		return faas
	},
	"kubeless": func() functions.FaaSDriver {
		faas, err := kubeless.New(&kubeless.Config{
			K8sConfig:       config.Global.Function.Kubeless.K8sConfig,
			FuncNamespace:   config.Global.Function.Kubeless.FuncNamespace,
			ImagePullSecret: config.Global.Function.Kubeless.ImagePullSecret,
		})
		if err != nil {
			log.Fatalf("Error starting Kubeless driver: %+v", err)
		}
		return faas
	},
	"noop": func() functions.FaaSDriver {
		faas, err := noop.New(&noop.Config{})
		if err != nil {
			log.Fatalf("Error starting noop driver: %+v", err)
		}
		return faas
	},
}

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

var debugFlags = struct {
	DebugEnabled bool `long:"debug" description:"Enable debugging messages"`
}{}

func configureFlags() []swag.CommandLineOptionsGroup {
	return []swag.CommandLineOptionsGroup{{
		ShortDescription: "Function manager Flags",
		LongDescription:  "",
		Options:          &functionmanager.FunctionManagerFlags,
	}, {
		ShortDescription: "Debug options",
		LongDescription:  "",
		Options:          &debugFlags,
	}}
}

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.FlatSwaggerJSON, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewFunctionManagerAPI(swaggerSpec)
	server := restapi.NewServer(api)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Function Manager"
	parser.LongDescription = "This is the API server for the Dispatch Function Manager service.\n"

	optsGroups := configureFlags()
	for _, optsGroup := range optsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	if debugFlags.DebugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	config.Global = config.LoadConfiguration(functionmanager.FunctionManagerFlags.Config)
	log.Debugln("config.Global:")
	log.Debugf("%+v", config.Global)

	registryAuth := config.Global.Registry.RegistryAuth
	if config.Global.Registry.RegistryAuth == "" {
		registryAuth = config.EmptyRegistryAuth
	}

	es, err := entitystore.NewFromBackend(
		entitystore.BackendConfig{
			Backend:  functionmanager.FunctionManagerFlags.DbBackend,
			Address:  functionmanager.FunctionManagerFlags.DbFile,
			Bucket:   functionmanager.FunctionManagerFlags.DbDatabase,
			Username: functionmanager.FunctionManagerFlags.DbUser,
			Password: functionmanager.FunctionManagerFlags.DbPassword,
		})
	if err != nil {
		log.Fatalln(err)
	}

	faas := drivers[config.Global.Function.Faas]()
	defer utils.Close(faas)

	c := &functionmanager.ControllerConfig{
		ResyncPeriod: time.Duration(config.Global.Function.ResyncPeriod) * time.Second,
	}

	secretsClient := client.NewSecretsClient(functionmanager.FunctionManagerFlags.SecretStore, client.AuthWithToken("cookie"), "")
	servicesClient := client.NewServicesClient(functionmanager.FunctionManagerFlags.ServiceManager, client.AuthWithToken("cookie"), "")

	r := runner.New(&runner.Config{
		Faas:            faas,
		Validator:       validator.New(),
		SecretInjector:  injectors.NewSecretInjector(secretsClient),
		ServiceInjector: injectors.NewServiceInjector(secretsClient, servicesClient),
	})

	var imageGetter functionmanager.ImageGetter
	imageGetter = client.NewImagesClient(functionmanager.FunctionManagerFlags.ImageManager, client.AuthWithToken("cookie"), "")
	if config.Global.Function.FileImageManager != "" {
		imageGetter = functionmanager.FileImageManagerClient()
	}

	dc, err := docker.NewEnvClient()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "could not get docker client"))
	}
	imageBuilder := functions.NewDockerImageBuilder(config.Global.Registry.RegistryURI, registryAuth, dc)

	controller := functionmanager.NewController(c, es, faas, r, imageGetter, imageBuilder)
	defer controller.Shutdown()
	controller.Start()

	handlers := functionmanager.NewHandlers(controller.Watcher(), es)
	handlers.ConfigureHandlers(api)

	healthChecker := func() error {
		// TODO: implement service-specific healthchecking
		return nil
	}

	tracer, tracingCloser, err := utils.CreateTracer("FunctionManager", functionmanager.FunctionManagerFlags.Tracer)
	if err != nil {
		log.Fatalf("Error creating a tracer: %+v", err)
	}
	defer tracingCloser.Close()
	opentracing.SetGlobalTracer(tracer)

	handler := alice.New(
		middleware.NewHealthCheckMW("", healthChecker),
		middleware.NewTracingMW(tracer),
	).Then(api.Serve(nil))

	server.SetHandler(handler)

	defer server.Shutdown()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
