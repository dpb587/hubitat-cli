package certificatecmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewUpdateCommand(cmdp *cmdflags.Persistent) *cobra.Command {
	var fCertificatePath, fPrivateKeyPath string
	var fForce bool
	var fSkipReboot bool

	var cmd = &cobra.Command{
		Use:   "update",
		Short: "For updating the certificate",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var err error
			var certificateBytes, privateKeyBytes []byte

			{ // certificate
				if len(fCertificatePath) == 0 {
					return fmt.Errorf("missing certificate flag")
				}

				certificateBytes, err = ioutil.ReadFile(fCertificatePath)
				if err != nil {
					return errors.Wrap(err, "reading certificate")
				}
			}

			{ // privateKey
				if len(fPrivateKeyPath) == 0 {
					return fmt.Errorf("missing private key flag")
				}

				privateKeyBytes, err = ioutil.ReadFile(fPrivateKeyPath)
				if err != nil {
					return errors.Wrap(err, "reading private key")
				}
			}

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			if fForce {
				cmdp.Logger.V(1).Info("forcing update")
			} else {
				cmdp.Logger.V(2).Info("comparing certificates")

				hubBaseURL := hubClient.BaseURL()

				if hubBaseURL.Scheme != "https" {
					cmdp.Logger.V(1).Info("unable to compare certificates for non-https url")
				} else {
					certPEM, _ := pem.Decode(certificateBytes)
					if certPEM == nil {
						return errors.New("decoding certificate: invalid data")
					}

					cert, err := x509.ParseCertificate(certPEM.Bytes)
					if err != nil {
						return fmt.Errorf("parsing certificate: %v", err)
					}

					tlsConfig, err := cmdp.TLSConfig()
					if err != nil {
						return errors.Wrap(err, "loading tls config")
					}

					var dialPort = hubBaseURL.Port()
					if dialPort == "" {
						dialPort = "443"
					}

					conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", hubBaseURL.Hostname(), dialPort), tlsConfig)
					if err != nil {
						return fmt.Errorf("tls dial: %v", err)
					}

					defer conn.Close()

					desiredSerialNumber := hex.EncodeToString(cert.SerialNumber.Bytes())
					actualSerialNumber := hex.EncodeToString(conn.ConnectionState().PeerCertificates[0].SerialNumber.Bytes())

					cmdp.Logger.V(1).Info("compared certificates", "desired", desiredSerialNumber, "actual", actualSerialNumber)

					if desiredSerialNumber == actualSerialNumber {
						cmdp.Logger.V(0).Info("skipping update (certificate already in use)", "sn", actualSerialNumber)

						return nil
					}
				}
			}

			err = hubClient.UpdateAdvancedCertificates(ctx, certificateBytes, privateKeyBytes)
			if err != nil {
				return errors.Wrap(err, "updating")
			}

			cmdp.Logger.V(0).Info("updated certificate")

			if fSkipReboot {
				return nil
			}

			err = hubClient.Reboot(ctx)
			if err != nil {
				return errors.Wrap(err, "rebooting")
			}

			cmdp.Logger.V(0).Info("requested reboot")

			return nil
		},
	}

	cmd.Flags().StringVar(&fCertificatePath, "certificate-path", "", "certificate file path (PEM format)")
	cmd.Flags().StringVar(&fPrivateKeyPath, "private-key-path", "", "private key file path (PEM format)")

	cmd.Flags().BoolVar(&fForce, "force", false, "force update (even if an update seems unnecessary)")
	cmd.Flags().BoolVar(&fSkipReboot, "skip-reboot", false, "skip reboot after update")

	return cmd
}
