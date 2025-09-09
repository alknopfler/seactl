# seactl (SUSE Edge Airgap Command Line Interface)

[![Coverage Status](https://coveralls.io/repos/github/alknopfler/seactl/badge.svg?branch=main)](https://coveralls.io/github/alknpofler/seactl?branch=main)


SUSE Edge Airgap Tool created to make the airgap process easier for SUSE Edge for telco deployments.

## Features

- Read the info from the airgap manifest file.
- Create a tarball for rke2 release tarball files (required to be used in capi airgap scenarios).
- Upload the helm-charts oci images defined in the release manifest to the private registry .
- Upload the containers images defined in the release manifest to the private registry.

## Requirements

- Helm 3 installed on the machine. You can install it using:

```shell
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

## Usage

Clone the repository and build the tool using the following command:

```shell
make compile
```

1. If your private registry is auth based, create your own registry auth file with the following format:

```txt
<username_bas64encoded>:<password_base64encoded>
```

for example you can generate both using
```
echo -n "myuser" | base64
echo -n "mypassword" | base64
```

2. If your private registry is using a self-signed certificate, create a CA certificate file and provide the path to the tool.

The following command can be used to generate the airgap artifacts

```bash
Usage:
seactl generate [flags]

Flags:
-h, --help                       help for generate
-i, --input string               Release manifest file
-k, --insecure                   Skip TLS verification in registry
-o, --output string              Output directory to store the tarball files
-a, --registry-authfile string   Registry Auth file with username:password base64 encoded
-c, --registry-cacert string     Registry CA Certificate file
-r, --registry-url string        Registry URL
-d, --dryrun                     Dry run mode, only print the actions without executing them
-m, --release-mode string        Release mode, can be 'factory' or 'production' (default "factory")
-v, --release-version string     Release version, e.g. 3.4.0 (X.Y.Z)
```

## Example of usage

```bash
seactl generate -v 3.4.0 -m factory -o /tmp/airgap -a registry-auth.txt -c /opt/certs/ca.crt -r myregistry:5000
```

```bash
seactl generate -v 3.4.0 -m production -o /tmp/airgap -a registry-auth.txt -r myregistry:5000 --insecure
```

```bash 
./seactl generate -v 3.4.0 -m production -o ./tmp/airgap -r localhost:3000 -d true
```