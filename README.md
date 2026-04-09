![logo](https://kvalitetsit.dk/wp-content/uploads/2024/06/Logo.png)

[![License](https://img.shields.io/badge/license-MIT-blue)](./LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-56.7%25-brightgreen)](https://codecov.io/gh/holsting/myra-cert-manager-webhook-rs)

# Myra Cert-Manager Webhook
A webhook component for handling **ACME DNS-01 challenges** via the Myra DNS provider. It integrates with [cert-manager](https://cert-manager.io/) to automate the issuance and renewal of TLS/SSL certificates.

---

## Features

- Automates creation of `_acme-challenge` DNS TXT records for domain validation
- Supports **wildcard certificates** (`*.example.com`)
- Handles **cleanup** of DNS records after validation
- Extensible to support other DNS providers besides Myra
- Containerized for easy deployment with Docker

---

## How It Works

This webhook acts as a bridge between **cert-manager**, the **ACME protocol**, and the **DNS provider**.  

1. A certificate request is issued via cert-manager.
2. The ACME server returns a **DNS-01 challenge** to verify domain ownership.
3. The webhook receives a `present` request and **creates the required TXT record** in DNS.
4. The ACME server queries DNS to verify the token.
5. After verification, the certificate is issued.
6. Finally, the webhook can handle a `cleanup` request to **remove the temporary TXT record**.

---

## Webhook
This webhook exposes endpoints that Cert-Manager calls to handle DNS-01 challenges when issuing certificates. It consists of two main steps: presenting the challenge and cleaning up after verification.

### Presenting the DNS Challenge

When Cert-Manager receives a DNS-01 challenge from the ACME server, it calls the webhook’s /present endpoint to create a TXT record for domain verification.

Example request:

```bash
curl localhost:8080/present \
  -H "Content-Type: application/json" \
  -d '{
        "domain": "_acme-challenge.example.com",
        "token": "<token>",
        "key_authorization": "<key>"
      }'
```

This is expected to create the following DNS record:
```txt
_acme-challenge.example.com TXT <token>
```

---

### Cleanup Phase

After the ACME server verifies the TXT record and the certificate is issued, Cert-Manager calls the /cleanup endpoint to remove the temporary TXT record.

Example request:

```bash
curl localhost:8080/cleanup \
  -H "Content-Type: application/json" \
  -d '{
        "domain": "_acme-challenge.example.com",
        "token": "<token>",
        "key_authorization": "<key>"
      }'
```

---

## Build & Deployment

### Build from Source
Build a standalone executable binary directly from the Go source code:
```bash
go build ./cmd/app
```
This outputs a binary named app in your current directory (or you can specify a path with -o).


### Build with Docker
Build a Docker image that packages the webhook and all its dependencies:

```bash
docker build -t kvalitetsit/myra-cert-manager-webhook:latest -f ./docker/Dockerfile .
```
This output a docker image tagged kvalitetsit/myra-cert-manager-webhook:latest. The resulting image can be deployed using [docker-compose](./docker/docker-compose.yaml).

### Build with Makefile
The Makefile provides a convenient shortcut to build the Docker image (and optionally manage other tasks like rendering manifests):

```bash
./make build
```
This outputs the same image as above. The Makefile automates setup and can be extended for local builds, testing, and deployment. You can combine it with other targets like make rendered-manifest.yaml for Kubernetes manifests.

## Test
The project aims to maintain a high and meaningful level of test coverage, focusing on critical paths and core functionality. The goal is to ensure reliability, reduce regressions, and support confident, maintainable development. A coverage report can be achieved by runnig the tests below with the `-cover` flag.

### Unit tests

Running the single command below should let you test the entire project recursivly.

```bash
go test ./...
```

### Integration
The integration tests can be executed with the following command. It is expected to start a mock of both the Myra API and a DNS. These mocks contribute to ensure the webhook integrates with the external interfaces.

```bash
TEST_ZONE_NAME=example.com. make test
```

> Ensure you are located in the root of the project when you run the commands above.
---

## Configuration

The webhook can be configured using **environment variables**. Below is a table of the supported options:

| Name                        | Description                                                         | Default               | Required | Notes |
|------------------------------|--------------------------------------------------------------------|---------------------|----------|-------|
| `MYRA_API_KEY`               | API key to authenticate with the Myra DNS provider                 | —                   | Yes      | Must have permissions to create/update TXT records |
| `MYRA_API_SECRET`            | API key to authenticate with the Myra DNS provider                 | —                   | Yes      | Must have permissions to create/update TXT records |
| `MYRA_API_TOKEN`             | API token to authenticate with the Myra DNS provider                  | —                   | Conditional      | Must have permissions to create/update TXT records required if the key/secret is not specified |
| `MYRA_API_URL`               | Base URL of the Myra DNS API                                        | `https://apiv2.myracloud.com` | No       | Override only if using a custom endpoint |
| `WEBHOOK_PORT`               | Port on which the webhook server listens                             | `8080`               | No       | Use a different port if 8080 is occupied |
| `LOG_LEVEL`                  | Logging verbosity level (`debug`, `info`, `warn`, `error`)          | `info`              | No       | `debug` recommended for troubleshooting |
| `ACME_CHALLENGE_TTL`         | TTL (time-to-live) for DNS TXT records created for validation       | `300` (seconds)     | No       | Adjust according to DNS propagation speed |
| `CERT_MANAGER_NAMESPACE`     | Kubernetes namespace where cert-manager operates                     | `default`           | No       | Useful for multi-namespace setups |
| `CUSTOM_DNS_CLIENT`          | Optional: path or name of a custom DNS client to replace Myra        | —                   | No       | Must implement the DNS client interface |

> **Tip:** For Kubernetes deployments, these variables can be set via a `Deployment` manifest as `env` entries.


---

## Usage

Once the webhook is deployed, create a `ClusterIssuer` that uses it as the DNS-01 solver:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          webhook:
            groupName: acme-myra.kvalitetsit.dk
            solverName: acme-myra.kvalitetsit.dk
```

---

## Extensibility

The webhook defaults to **Myra DNS**, but can be **extended with a custom DNS client**. This allows integration with other DNS providers or internal systems by implementing the same simple client interface below.

The DNS client interface:
```go
type Client[T any] interface {
	OnDelete(record T) (T, error)
	OnAdd(record T) (T, error)
}
```

---

## License

This project is licensed under the [MIT License](./LICENSE).

---

## References

- [ACME Protocol (RFC 8555)](https://datatracker.ietf.org/doc/html/rfc8555)
- [cert-manager](https://cert-manager.io/)
- [Myra Secure DNS](https://www.myrasecurity.com/en/saasp/application-security/secure-dns/)

