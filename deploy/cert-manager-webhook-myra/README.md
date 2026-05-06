# cert-manager-webhook-myra

Helm chart for deploying the [Myra](https://www.myrasecurity.com/en/) cert-manager ACME DNS-01 webhook solver.

## Prerequisites

- Kubernetes 1.21+
- Helm 3.0+
- [cert-manager](https://cert-manager.io/) v1.19+ installed in the cluster
- A Myra API key and secret with permissions to create/delete TXT records

## Installation

```bash
helm repo add kvalitetsit https://raw.githubusercontent.com/KvalitetsIT/helm-repo/master/
helm repo update

helm install cert-manager-webhook-myra kvalitetsit/cert-manager-webhook-myra \
  --namespace cert-manager \
  --set myra.apiKeySecretName=myra-api-credentials
```

The chart expects a Kubernetes secret with the Myra API credentials:

```bash
kubectl create secret generic myra-api-credentials \
  --namespace cert-manager \
  --from-literal=api-key=<your-api-key> \
  --from-literal=api-secret=<your-api-secret>
```

## Usage

After installation, create a `ClusterIssuer` that uses the webhook as its DNS-01 solver:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: your@email.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          webhook:
            groupName: acme-myra.kvalitetsit.dk
            solverName: acme-myra.kvalitetsit.dk
```

## Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `groupName` | Unique group name for the webhook API group | `acme-myra.kvalitetsit.dk` |
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `kvalitetsit/cert-manager-webhook-myra` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `myra.apiKeySecretName` | Name of the secret containing the Myra API credentials | `myra-api-credentials` |
| `certManager.namespace` | Namespace where cert-manager is installed | `cert-manager` |
| `certManager.serviceAccountName` | Service account name of cert-manager | `cert-manager` |
| `resources.limits.memory` | Memory limit | `100Mi` |
| `resources.requests.cpu` | CPU request | `5m` |
| `resources.requests.memory` | Memory request | `30Mi` |
| `podSecurityContext` | Pod-level security context | `runAsNonRoot: true`, `runAsUser/Group/fsGroup: 65532` |
| `securityContext` | Container-level security context | `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, capabilities `ALL` dropped |
| `networkPolicy.enabled` | Enable NetworkPolicy | `true` |
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity rules | `{}` |
