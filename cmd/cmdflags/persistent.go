package cmdflags

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/dpb587/hubitat-cli/hub"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/publicsuffix"
)

type Persistent struct {
	Stdout io.Writer
	Stderr io.Writer

	Logger logr.Logger

	fVerbosity int

	fHubURL      string
	fHubCAPath   string
	fHubInsecure bool
	fHubUsername string
	fHubPassword string

	hubClient *hub.Client
	tlsConfig *tls.Config
}

func NewPersistent(cmd *cobra.Command) *Persistent {
	p := &Persistent{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.PersistentPreRunE = p.bind(strings.Fields(cmd.Use)[0], cmd.PersistentFlags(), cmd.PersistentPreRunE)

	return p
}

func (p *Persistent) bind(name string, flags *pflag.FlagSet, preRunE func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	flags.IntVarP(&p.fVerbosity, "verbose", "v", 0, "verbosity level")

	flags.StringVar(&p.fHubURL, "hub-url", "", "URL for accessing the hub (e.g. http://192.0.2.100; $HUBITAT_URL)")
	flags.StringVar(&p.fHubCAPath, "hub-ca-path", "", "custom certificate authorities to trust ($HUBITAT_CA_PATH)")
	flags.BoolVar(&p.fHubInsecure, "hub-insecure", false, "disable TLS verifications ($HUBITAT_INSECURE)")
	flags.StringVar(&p.fHubUsername, "hub-username", "", "username for login ($HUBITAT_USERNAME)")
	flags.StringVar(&p.fHubPassword, "hub-password", "", "password for login ($HUBITAT_PASSWORD)")

	return func(cmd *cobra.Command, args []string) error {
		if preRunE != nil {
			if err := preRunE(cmd, args); err != nil {
				return err
			}
		}

		for k, ev := range map[string]string{
			"hub-url":      "HUBITAT_URL",
			"hub-ca-path":  "HUBITAT_CA_PATH",
			"hub-insecure": "HUBITAT_INSECURE",
			"hub-username": "HUBITAT_USERNAME",
			"hub-password": "HUBITAT_PASSWORD",
		} {
			pf := flags.Lookup(k)
			if pf.Changed {
				continue
			} else if v, ok := os.LookupEnv(ev); ok {
				if err := pf.Value.Set(v); err != nil {
					return errors.Wrapf(err, "setting %s from env", k)
				}
			}
		}

		stdr.SetVerbosity(p.Verbosity())
		p.Logger = stdr.NewWithOptions(log.New(p.Stderr, "", log.LstdFlags), stdr.Options{}).WithName(name)
		p.Logger.V(1).Info("runtime", "os", runtime.GOOS, "arch", runtime.GOARCH)
		p.Logger.V(1).Info("version", "name", VersionName, "commit", VersionCommit, "built", VersionBuilt)

		return nil
	}
}

func (p *Persistent) HubClient() (*hub.Client, error) {
	if p.hubClient == nil {
		if len(p.fHubURL) == 0 {
			return nil, errors.New("missing hub url")
		}

		hubURL, err := url.Parse(p.fHubURL)
		if err != nil {
			return nil, errors.Wrap(err, "parsing hub url")
		}

		tlsConfig, err := p.TLSConfig()
		if err != nil {
			return nil, errors.Wrap(err, "loading tls config")
		}

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return nil, errors.Wrap(err, "creating cookiejar")
		}

		hubClient := hub.NewClient(
			p.Logger,
			&http.Client{
				Jar:     jar,
				Timeout: 60 * time.Second,
				Transport: &http.Transport{
					Dial: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
					}).Dial,
					IdleConnTimeout:       30 * time.Second,
					TLSClientConfig:       tlsConfig,
					TLSHandshakeTimeout:   30 * time.Second,
					ResponseHeaderTimeout: 60 * time.Second,
					ExpectContinueTimeout: 15 * time.Second,
				},
			},
			hubURL,
		)

		if len(p.fHubUsername) > 0 || len(p.fHubPassword) > 0 {
			err := hubClient.Login(context.TODO(), p.fHubUsername, p.fHubPassword)
			if err != nil {
				return nil, errors.Wrap(err, "logging in")
			}
		}

		p.hubClient = hubClient
	}

	return p.hubClient, nil
}

func (p *Persistent) TLSConfig() (*tls.Config, error) {
	if p.tlsConfig == nil {
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		if p.fHubCAPath != "" {
			rootCAs = x509.NewCertPool()

			certs, err := ioutil.ReadFile(p.fHubCAPath)
			if err != nil {
				return nil, errors.Wrap(err, "reading custom ca certs")
			}

			ok := rootCAs.AppendCertsFromPEM(certs)
			if !ok {
				return nil, fmt.Errorf("unable to append custom ca certs")
			}
		}

		if p.fHubInsecure {
			p.Logger.V(0).Info("tls verifications have been disabled")
		}

		p.tlsConfig = &tls.Config{
			InsecureSkipVerify: p.fHubInsecure,
			RootCAs:            rootCAs,
		}
	}

	return p.tlsConfig, nil
}

func (p *Persistent) Verbosity() int {
	return p.fVerbosity
}
