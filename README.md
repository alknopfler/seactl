# seactl
SUSE Edge Airgap Tool created to make the airgap process easier for SUSE Edge deployments.

## Features

- Read the info from the release manifest file.
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
```

## Example

```bash
seactl generate -i release-manifest.yaml -o /tmp/airgap -a registry-auth.txt -c /opt/certs/ca.crt -r myregistry:5000
```

```bash
seactl generate -i release-manifest.yaml -o /tmp/airgap -a registry-auth.txt -r myregistry:5000 --insecure
```