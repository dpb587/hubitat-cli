# hubitat-cli

An unofficial CLI for some [Hubitat](https://hubitat.com/) commands.

Hubitat doesn't really have a formal API, so this is a little hacky and currently only supports a couple management commands I rely on in automation (namely downloading backups and rotating HTTPS certificates).

## Usage

Download and install the binary for your Linux, macOS, or Windows environment from the [releases page](https://github.com/dpb587/hubitat-cli/releases). At a minimum, you'll need to configure the IP or hostname of your Hubitat node.

```bash
$ hubitat-cli --hub-url='http://192.0.2.100' ...
```

If your hub requires login ([learn more](https://docs.hubitat.com/index.php?title=Hub_Login_Security)), you will also want to configure your username and password.

```bash
$ hubitat-cli \
  --hub-url='http://192.0.2.100' \
  --hub-username='janedoe' \
  --hub-password='secretpassword' \
  ...
```

Environment variables may be used for many of the global flags to help with automation or avoid lengthy commands.

```bash
$ export HUBITAT_URL=http://192.0.2.100
$ export HUBITAT_USERNAME='janedoe'
$ export HUBITAT_PASSWORD='secretpassword'
$ hubitat-cli ...
```

Use the `--help` flag at any point to view more details on flags, arguments, and any available subcommands. For example, to download the latest backup you might use the following.

```bash
$ hubitat-cli backup download --output latest.lzf
```

Or, you might want to list the backup files with the `list` subcommand instead.

```bash
$ hubitat-cli backup list
```

### Advanced

A minimal image is available for use from [Docker](https://www.docker.com/), [Kubernetes](https://kubernetes.io/), or any other container runtime.

```bash
$ docker run ghcr.io/dpb587/hubitat-cli/main -- version
```

If Hubitat has HTTPS enabled, the CLI will reject any certificates your system does not already trust. To trust a custom CA, configure the CA certificate file path with the `--hub-ca-cert=` flag (or `HUBITAT_CA_CERT` environment variable). See the `--help` option for additional connection options and methods for disabling secure connections.

## License

[MIT License](LICENSE)
