///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/vmware/dispatch/pkg/dispatchcli/i18n"
)

var (
	dispatchHelmRepositoryURL = "https://s3-us-west-2.amazonaws.com/dispatch-charts"
	servicesAvailable         = []string{"certs", "ingress", "postgres", "api-gateway", "dispatch"}
	servicesEnabled           map[string]bool
)

type chartConfig struct {
	Chart     string `json:"chart,omitempty" validate:"required"`
	Namespace string `json:"namespace,omitempty" validate:"required,hostname"`
	Release   string `json:"release,omitempty" validate:"required"`
	Repo      string `json:"repo,omitempty" validate:"omitempty,uri"`
	Version   string `json:"version,omitempty" validate:"omitempty"`
}

type ingressConfig struct {
	Chart       *chartConfig `json:"chart,omitempty" validate:"required"`
	ServiceType string       `json:"serviceType,omitempty" validate:"required,eq=NodePort|eq=LoadBalancer|eq=ClusterIP"`
}

type postgresConfig struct {
	Chart       *chartConfig `json:"chart,omitempty" validate:"required"`
	Database    string       `json:"database,omitempty" validate:"required"`
	Username    string       `json:"username,omitempty" validate:"required"`
	Password    string       `json:"password,omitempty" validate:"required"`
	Host        string       `json:"host,omitempty" validate:"required,hostname"`
	Port        int          `json:"port,omitempty" validate:"required"`
	Persistence bool         `json:"persistence,omitempty" validate:"omitempty"`
}

type tlsConfig struct {
	CertFile   string `json:"certFile,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	SecretName string `json:"secretName,omitempty" validate:"required"`
}

type apiGatewayConfig struct {
	Chart       *chartConfig `json:"chart,omitempty" validate:"required"`
	ServiceType string       `json:"serviceType,omitempty" validate:"required,eq=NodePort|eq=LoadBalancer|eq=ClusterIP"`
	Database    string       `json:"database,omitempty" validate:"required"`
	Hostname    string       `json:"hostname,omitempty" validate:"required,hostname"`
	TLS         *tlsConfig   `json:"tls,omitempty" validate:"required"`
}

type openfaasConfig struct {
	Chart         *chartConfig `json:"chart,omitempty" validate:"required"`
	ExposeService bool         `json:"exposeService,omitempty" validate:"omitempty"`
}

type riffConfig struct {
	Chart *chartConfig `json:"chart,omitempty" validate:"required"`
}

type imageConfig struct {
	Host string `json:"host,omitempty" validate:"required,hostname"`
	Tag  string `json:"tag,omitempty"  validate:"required"`
}
type oauth2ProxyConfig struct {
	ClientID     string `json:"clientID,omitempty" validate:"required"`
	ClientSecret string `json:"clientSecret,omitempty" validate:"required"`
	CookieSecret string `json:"cookieSecret,omitmepty" validate:"omitempty"`
}
type dockerRegistryConfig struct {
	Host     string `json:"host,omitempty" validate:"required,hostname"`
	Password string `json:"password,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Username string `json:"username,omitempty" validate:"required"`
}
type dispatchInstallConfig struct {
	FaaSBackend    string               `json:"faasbackend,omitempty" validate:"required,eq=openfaas|eq=riff"`
	Chart          *chartConfig         `json:"chart,omitempty" validate:"required"`
	Hostname       string               `json:"hostname,omitempty" validate:"required,hostname"`
	Port           int                  `json:"port,omitempty" validate:"required"`
	Organization   string               `json:"organization,omitempty" validate:"required"`
	Image          *imageConfig         `json:"image,omitempty" validate:"required"`
	Debug          bool                 `json:"debug,omitempty" validate:"omitempty"`
	Trace          bool                 `json:"trace,omitempty" validate:"omitempty"`
	Database       string               `json:"database,omitempty" validate:"required,eq=postgres"`
	PersistData    bool                 `json:"persistData,omitempty" validate:"omitempty"`
	DockerRegistry dockerRegistryConfig `json:"dockerRegistry,omitempty" validate:"required"`
	OAuth2Proxy    *oauth2ProxyConfig   `json:"oauth2Proxy,omitempty" validate:"required"`
	TLS            *tlsConfig           `json:"tls,omitempty" validate:"required"`
}

type installConfig struct {
	HelmRepositoryURL string                 `json:"helmRepositoryUrl,omitempty" validate:"required,uri"`
	Ingress           *ingressConfig         `json:"ingress,omitempty" validate:"required"`
	PostgresConfig    *postgresConfig        `json:"postgresql,omitempty" validate:"required"`
	APIGateway        *apiGatewayConfig      `json:"apiGateway,omitempty" validate:"required"`
	OpenFaas          *openfaasConfig        `json:"openfaas,omitempty" validate:"required"`
	Riff              *riffConfig            `json:"riff,omitempty" validate:"required"`
	DispatchConfig    *dispatchInstallConfig `json:"dispatch,omitempty" validate:"required"`
}

var (
	installLong = `Install the Dispatch framework.`

	installExample    = i18n.T(``)
	installConfigFile = i18n.T(``)
	installServices   []string
	installChartsDir  = i18n.T(``)
	installChartsRepo = i18n.T(``)
	installDryRun     = false
	installDebug      = false
	configDest        = i18n.T(``)
)

// NewCmdInstall creates a command object for the generic "get" action, which
// retrieves one or more resources from a server.
func NewCmdInstall(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [flags]",
		Short:   i18n.T("Install some or all of dispatch"),
		Long:    installLong,
		Example: installExample,
		Run: func(cmd *cobra.Command, args []string) {
			if installConfigFile == "" {
				runHelp(cmd, args)
				return
			}
			err := runInstall(out, errOut, cmd, args)
			CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&installConfigFile, "file", "f", "", "Path to YAML file")
	cmd.Flags().StringArrayVarP(&installServices, "service", "s", []string{}, "Service to install (defaults to all). Add 'no-' prefix to service name to disable it.")
	cmd.Flags().BoolVar(&installDryRun, "dry-run", false, "Do a dry run, but don't install anything")
	cmd.Flags().BoolVar(&installDebug, "debug", false, "Extra debug output")
	cmd.Flags().StringVar(&installChartsRepo, "charts-repo", "dispatch", "Helm Chart Repo used")
	cmd.Flags().StringVar(&installChartsDir, "charts-dir", "dispatch", "File path to local charts (for chart development)")
	cmd.Flags().StringVarP(&configDest, "destination", "d", "~/.dispatch", "Destination of the CLI configuration")
	return cmd
}

func installCert(out, errOut io.Writer, configDir, namespace, domain string, tls *tlsConfig) error {
	var key, cert string
	if tls.CertFile != "" {
		if tls.PrivateKey == "" {
			return errors.New("Error installing certificate: missing private key for the tls cert.")
		}
		key = tls.PrivateKey
		cert = tls.CertFile
	} else {
		// make a new key and cert.
		subject := fmt.Sprintf("/CN=%s/O=%s", domain, domain)
		key = path.Join(configDir, fmt.Sprintf("%s.key", domain))
		cert = path.Join(configDir, fmt.Sprintf("%s.crt", domain))
		var err error
		// If cert and key exist, reuse them
		if _, err = os.Stat(key); os.IsNotExist(err) {
			if _, err = os.Stat(cert); os.IsNotExist(err) {
				openssl := exec.Command(
					"openssl", "req", "-x509", "-nodes", "-days", "365", "-newkey", "rsa:2048",
					"-keyout", key,
					"-out", cert,
					"-subj", subject)
				opensslOut, err := openssl.CombinedOutput()
				if err != nil {
					return errors.Wrapf(err, string(opensslOut))
				}
			}
		}
	}
	kubectl := exec.Command(
		"kubectl", "delete", "secret", "tls", tls.SecretName, "-n", namespace)
	kubectlOut, err := kubectl.CombinedOutput()
	if err != nil {
		if !strings.Contains(string(kubectlOut), "NotFound") {
			return errors.Wrapf(err, string(kubectlOut))
		}
	}
	kubectl = exec.Command(
		"kubectl", "create", "namespace", namespace)
	kubectlOut, err = kubectl.CombinedOutput()
	if err != nil {
		if !strings.Contains(string(kubectlOut), "AlreadyExists") {
			return errors.Wrapf(err, string(kubectlOut))
		}
	}
	kubectl = exec.Command(
		"kubectl", "create", "secret", "tls", tls.SecretName, "-n", namespace, "--key", key, "--cert", cert)
	kubectlOut, err = kubectl.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, string(kubectlOut))
	}
	return nil
}

func helmRepoUpdate(out, errOut io.Writer, name, repoURL string) error {
	helm := exec.Command(
		"helm", "repo", "add", name, repoURL)
	helmOut, err := helm.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, string(helmOut))
	}
	helm = exec.Command("helm", "repo", "update")
	helmOut, err = helm.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, string(helmOut))
	}
	return nil
}

func helmDepUp(out, errOut io.Writer, chart string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Error getting current working directory")
	}
	err = os.Chdir(chart)
	if err != nil {
		return errors.Wrap(err, "Error changing directory")
	}
	helm := exec.Command("helm", "dep", "up")
	helmOut, err := helm.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, string(helmOut))
	}
	return os.Chdir(cwd)
}

func helmInstall(out, errOut io.Writer, meta *chartConfig, options map[string]string) error {

	repo := ""
	chart := meta.Chart
	if meta.Repo != "" {
		// if user specified a repo, use that
		repo = meta.Repo
	} else if installChartsDir == "dispatch" {
		// use dispatch chart repo
		repo = dispatchHelmRepositoryURL
	} else {
		// use the local charts
		chart = path.Join(installChartsDir, meta.Chart)
	}

	args := []string{"upgrade", "-i", meta.Release, chart, "--namespace", meta.Namespace}
	for k, v := range options {
		args = append(args, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	if repo != "" {
		args = append(args, "--repo", repo)
	}
	if meta.Version != "" {
		args = append(args, "--version", meta.Version)
	}
	if installDebug {
		args = append(args, "--debug")
	}
	args = append(args, "--wait")
	if installDryRun {
		args = append(args, "--dry-run")
	}
	if installDebug {
		fmt.Fprint(out, "debug: helm")
		for _, a := range args {
			fmt.Fprintf(out, " %s", a)
		}
		fmt.Fprint(out, "\n")
	}

	helm := exec.Command("helm", args...)
	helmOut, err := helm.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, string(helmOut))
	}
	if installDebug {
		fmt.Fprintln(out, string(helmOut))
	}
	return nil
}

func writeConfig(out, errOut io.Writer, configDir string, config *installConfig) error {
	dispatchConfig.Organization = config.DispatchConfig.Organization
	dispatchConfig.Host = config.DispatchConfig.Hostname
	dispatchConfig.Port = config.DispatchConfig.Port
	b, err := json.MarshalIndent(dispatchConfig, "", "    ")
	if err != nil {
		return err
	}
	if installDryRun {
		fmt.Fprintf(out, "Copy the following to your %s/config.json\n", configDir)
		fmt.Fprintln(out, string(b))
	} else {
		configPath := path.Join(configDir, "config.json")
		fmt.Fprintf(out, "Config file written to: %s\n", configPath)
		return ioutil.WriteFile(configPath, b, 0644)
	}
	return nil
}

func readConfig(out, errOut io.Writer, file string) (*installConfig, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", file)
	}
	config := installConfig{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding yaml file %s", file)
	}

	defaultInstallConfig := installConfig{}
	err = yaml.Unmarshal([]byte(defaultInstallConfigYaml), &defaultInstallConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding default install config yaml file")
	}
	if installDebug {
		b, _ := json.MarshalIndent(config, "", "    ")
		fmt.Fprintln(out, "Config before merge")
		fmt.Fprintln(out, string(b))
	}
	err = mergo.Merge(&config, defaultInstallConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "error merging default values")
	}
	return &config, nil
}

func selectServices(out io.Writer, config *installConfig) error {
	for _, service := range servicesAvailable {
		servicesEnabled[service] = true
	}

	servicesEnabled[config.DispatchConfig.FaaSBackend] = true

	// Most used combination - all services enabled and default backend
	if len(installServices) == 0 || installServices[0] == "all" {
		return nil
	}

	// We have two modes: whitelisting or blacklisting. Adding "no-" prefix to the service name
	// enters blacklist mode. Modes cannot be mixed.
	var whitelistMode, blacklistMode bool

	for _, service := range installServices {
		if strings.HasPrefix(service, "no-") {
			blacklistMode = true
			service := strings.TrimPrefix(service, "no-")
			if _, ok := servicesEnabled[service]; !ok {
				return fmt.Errorf("unknown service: %s", service)
			}

			if whitelistMode {
				return errors.New("can either whitelist or blacklist services, not both.")
			}

			servicesEnabled[service] = false
		} else {
			// whitelist mode
			if !whitelistMode {
				// we entered whitelist mode, during first encounter we need to reset enabled services
				for service := range servicesEnabled {
					servicesEnabled[service] = false
				}
				whitelistMode = true
			}

			if _, ok := servicesEnabled[service]; !ok {
				return fmt.Errorf("unknown service: %s", service)
			}

			if blacklistMode {
				return errors.New("can either whitelist or blacklist services, not both.")
			}
			servicesEnabled[service] = true
		}
	}
	if installDebug {
		fmt.Fprint(out, "\nServices to be installed:\n")
		for service, enabled := range servicesEnabled {
			if enabled {
				fmt.Fprintf(out, "* %s\n", service)
			}
		}
	}
	return nil
}

func installService(service string) bool {
	return servicesEnabled[service]
}

func getK8sServiceNodePort(service, namespace string, https bool) (int, error) {

	standardPort := 80
	if https {
		standardPort = 443
	}

	fmt.Printf("get nodePort for service %s from namespace %s", service, namespace)
	kubectl := exec.Command(
		"kubectl", "get", "svc", service, "-n", namespace,
		"-o", fmt.Sprintf("jsonpath={.spec.ports[?(@.port==%d)].nodePort}", standardPort))

	kubectlOut, err := kubectl.CombinedOutput()
	if err != nil {
		return -1, errors.Wrapf(err, string(kubectlOut))
	}
	nodePort, err := strconv.Atoi(string(kubectlOut))
	if err != nil {
		return -1, errors.Wrapf(err, "error fetching node port")
	}
	return nodePort, nil
}

func runInstall(out, errOut io.Writer, cmd *cobra.Command, args []string) error {

	config, err := readConfig(out, errOut, installConfigFile)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return errors.Wrapf(err, "configuration error")
	}

	if config.HelmRepositoryURL != "" {
		dispatchHelmRepositoryURL = config.HelmRepositoryURL
	}

	if installDebug {
		b, _ := json.MarshalIndent(config, "", "    ")
		fmt.Fprintln(out, string(b))
	}

	configDir, err := homedir.Expand(configDest)
	if !installDryRun {
		_, err = os.Stat(configDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(configDir, 0755)
			if err != nil {
				return errors.Wrapf(err, "error creating config destination directory")
			}
		}
	}

	selectServices(out, config)

	if installService("certs") || !installDryRun {
		err = installCert(out, errOut, configDir, config.DispatchConfig.Chart.Namespace, config.DispatchConfig.Hostname, config.DispatchConfig.TLS)
		if err != nil {
			return errors.Wrapf(err, "error installing cert for %s", config.DispatchConfig.Hostname)
		}
		err = installCert(out, errOut, configDir, config.APIGateway.Chart.Namespace, config.APIGateway.Hostname, config.APIGateway.TLS)
		if err != nil {
			return errors.Wrapf(err, "error installing  cert for %s", config.APIGateway.Hostname)
		}
	}
	if installChartsDir == "dispatch" {
		err = helmRepoUpdate(out, errOut, installChartsDir, config.HelmRepositoryURL)
		if err != nil {
			return errors.Wrapf(err, "error updating helm")
		}
	}

	if installService("ingress") {
		ingressOpts := map[string]string{
			"controller.service.type": config.Ingress.ServiceType,
		}
		err = helmInstall(out, errOut, config.Ingress.Chart, ingressOpts)
		if err != nil {
			return errors.Wrapf(err, "error installing nginx-ingress chart")
		}
		if config.Ingress.ServiceType == "NodePort" {
			service := fmt.Sprintf("%s-nginx-ingress-controller", config.Ingress.Chart.Release)
			config.DispatchConfig.Port, err = getK8sServiceNodePort(service, config.Ingress.Chart.Namespace, true)
			if err != nil {
				return err
			}
		}
	}

	if installService("postgres") {
		postgresOpts := map[string]string{
			"postgresDatabase":    config.PostgresConfig.Database,
			"postgresUser":        config.PostgresConfig.Username,
			"postgresPassword":    config.PostgresConfig.Password,
			"postgresHost":        config.PostgresConfig.Host,
			"postgresPort":        fmt.Sprintf("%d", config.PostgresConfig.Port),
			"persistence.enabled": strconv.FormatBool(config.PostgresConfig.Persistence),
		}
		err = helmInstall(out, errOut, config.PostgresConfig.Chart, postgresOpts)
		if err != nil {
			return errors.Wrapf(err, "error installing postgres chart")
		}
	}

	if installService("api-gateway") {
		kongOpts := map[string]string{
			"services.proxyService.type":  config.APIGateway.ServiceType,
			"database":                    config.APIGateway.Database,
			"postgresql.postgresDatabase": config.PostgresConfig.Database,
			"postgresql.postgresUser":     config.PostgresConfig.Username,
			"postgresql.postgresPassword": config.PostgresConfig.Password,
			"postgresql.postgresHost":     config.PostgresConfig.Host,
			"postgresql.postgresPort":     fmt.Sprintf("%d", config.PostgresConfig.Port),
		}
		err = helmInstall(out, errOut, config.APIGateway.Chart, kongOpts)
		if err != nil {
			return errors.Wrapf(err, "error installing kong chart")
		}

		if config.APIGateway.ServiceType == "NodePort" {

			service := fmt.Sprintf("%s-kongproxy", config.APIGateway.Chart.Release)
			httpsPort, err := getK8sServiceNodePort(service, config.APIGateway.Chart.Namespace, true)
			if err != nil {
				return err
			}
			httpPort, err := getK8sServiceNodePort(service, config.APIGateway.Chart.Namespace, false)
			if err != nil {
				return err
			}

			os.Setenv("API_GATEWAY_HTTPS_PORT", strconv.Itoa(httpsPort))
			os.Setenv("API_GATEWAY_HTTP_PORT", strconv.Itoa(httpPort))
		}
	}

	if installService("openfaas") {
		openFaasOpts := map[string]string{"exposeServices": "false"}
		err = helmInstall(out, errOut, config.OpenFaas.Chart, openFaasOpts)
		if err != nil {
			return errors.Wrapf(err, "error installing Riff chart")
		}
	}

	if installService("riff") {
		riffOpts := map[string]string{"exposeServices": "false"}
		err = helmInstall(out, errOut, config.Riff.Chart, riffOpts)
		if err != nil {
			return errors.Wrapf(err, "error installing Riff chart")
		}
	}

	if installService("dispatch") {
		chart := path.Join(installChartsDir, "dispatch")
		if installChartsDir != "dispatch" {
			err = helmDepUp(out, errOut, chart)
			if err != nil {
				return errors.Wrap(err, "error updating chart dependencies")
			}
		}

		// Resets the cookie every deployment if not specified
		if config.DispatchConfig.OAuth2Proxy.CookieSecret == "" {
			cookie := make([]byte, 16)
			_, err := rand.Read(cookie)
			if err != nil {
				return errors.Wrap(err, "error creating cookie secret")
			}
			config.DispatchConfig.OAuth2Proxy.CookieSecret = base64.StdEncoding.EncodeToString(cookie)
		}

		registryAuth := fmt.Sprintf(
			`{"username":"%s","password":"%s","email":"%s"}`,
			config.DispatchConfig.DockerRegistry.Username,
			config.DispatchConfig.DockerRegistry.Password,
			config.DispatchConfig.DockerRegistry.Email)

		registryAuthEncoded := base64.StdEncoding.EncodeToString([]byte(registryAuth))
		dispatchOpts := map[string]string{
			"global.host":                                  config.DispatchConfig.Hostname,
			"global.port":                                  strconv.Itoa(config.DispatchConfig.Port),
			"global.debug":                                 strconv.FormatBool(config.DispatchConfig.Debug),
			"global.trace":                                 strconv.FormatBool(config.DispatchConfig.Trace),
			"global.data.persist":                          strconv.FormatBool(config.DispatchConfig.PersistData),
			"function-manager.faas.openfaas.registryAuth":  registryAuthEncoded,
			"function-manager.faas.openfaas.imageRegistry": config.DispatchConfig.DockerRegistry.Host,
			"oauth2-proxy.app.clientID":                    config.DispatchConfig.OAuth2Proxy.ClientID,
			"oauth2-proxy.app.clientSecret":                config.DispatchConfig.OAuth2Proxy.ClientSecret,
			"oauth2-proxy.app.cookieSecret":                config.DispatchConfig.OAuth2Proxy.CookieSecret,
			"global.db.backend":                            config.DispatchConfig.Database,
			"global.db.host":                               config.PostgresConfig.Host,
			"global.db.port":                               fmt.Sprintf("%d", config.PostgresConfig.Port),
			"global.db.user":                               config.PostgresConfig.Username,
			"global.db.password":                           config.PostgresConfig.Password,
			"global.db.release":                            config.PostgresConfig.Chart.Release,
			"global.db.namespace":                          config.PostgresConfig.Chart.Namespace,
		}

		// If unset values default to chart values
		if config.DispatchConfig != nil && config.DispatchConfig.Image != nil {
			if config.DispatchConfig.Image.Host != "" {
				dispatchOpts["global.image.host"] = config.DispatchConfig.Image.Host
			}
			if config.DispatchConfig.Image.Tag != "" {
				dispatchOpts["global.image.tag"] = config.DispatchConfig.Image.Tag
			}
		}
		if installDebug {
			for k, v := range dispatchOpts {
				fmt.Fprintf(out, "%v: %v\n", k, v)
			}
		}
		err = helmInstall(out, errOut, config.DispatchConfig.Chart, dispatchOpts)
		if err != nil {
			return errors.Wrapf(err, "Error installing dispatch chart")
		}
	}
	err = writeConfig(out, errOut, configDir, config)
	return err
}
