---
authors: Diego Rodriguez (diego@rodilla.me)
state: draft
---

# RFD 0 - SRE Server Design

## What
Initial design for SRE server HTTP API server that provides interactions with a
Kubernetes cluster.

## Why

As part of the SRE's job, they are required to interact with Kubernetes clusters
on regular basis. However, existing tools such as `kubectl` or using the
Kubernetes API directly, can pose considerable friction with the overwhelming
amount options that are available in those tools. Also, onboarding new SREs with
those tools can be non-trivial, requring them to spend weeks learning the
terminology and how to actually use them. To reduce the friction of interaction
with Kubernetes clusters and ramp up time, we therefore propose a new service
`sre-server` to help ease these pains. The `sre-server` will provide convenience
(and safety) on top of the existing Kuberenets API to simplify interactions and
thus save the SRE's time while exponentially increasing their productivity.

## Details

### Scope
After collecting feedback from the SREs, it turns out that they spend a large
chunk of their time working with Kubernetes deployments. Hence, for the first
pass, we will focus primarily on providing niceities over the Kubernetes
deployments. Depending on the success of this project, we may decide to further
expand the scope to other Kubernetes resources as well.


### API
All responses and request bodies shall be formated in JSON.


In general, all non 2XX/3XX statuses shall return a response as follows:

```json
{
  "message" : "Message explaining the status."
}
```

We limit the scope of the RBAC permissions on the `Deployments` resource to only
what is needed; we try to follow the principle of least privilege for enhanced
security.


#### **`GET`** `/deployments[/{namespace}]`
Lists all of the deployments of the Kubernetes cluster by namespace

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The namespace to limit the scope of the deployments to list. Omitting it the namespace will retrieve all the deployments of cluster.  | `false` |


*Response*
```json
[
  {
    "namespace": "default",
    "deployments": [
      "foo",
      "yeet",
      ...
    ]
  },
  ...
]
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully retrieved the list of deployments |
| 404 | The provided namespace does not exist |




#### **`GET`** `/deployments/{namespace}/{deployment}/replicas`
Retrieves the deployment replica count by name in the provided namespace.

For this particular endpoint, to prevent overwhelming the Kubernetes cluster API
with too many of these requests (impacting its performance), the replica count
shall be cached. For the initial implementation, it shall use an in-memory cache
implemented as a map for simplicity. Due to Go maps being [unsafe
for concurrent use](https://go.dev/doc/faq#atomic_maps), we'll need to synchronize
the read/write access (which may impact performance). Eventually, for better
efficiency and performance across multiple replicas, we'd want to use something
like Redis to centralize the cache.

In addition to caching the replica count, when a new entry is added, we will add
a deployment watcher to check for any updates to the replica count to keep the
cache up to date. In the event that the deployment no longer exists, the
associated cache entry shall eventually be removed by the watcher.

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The deployment namespace. | `true` |
| `deployment`   | The name of the deployment. | `true` |


*Response*
```json
{
  "namespace": "default",
  "deployment": "cool",
  "replicas": 7
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully retrieved the replica count for the deployment |
| 404 | The provided namespace or deployment does not exist. |




#### **`PUT`** `/deployments/{namespace}/{deployment}/replicas`
Update the replica count for the given deployment.

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The deployment namespace. | `true` |
| `deployment`   | The name of the deployment. | `true` |

| Body Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `replicas`    | The new number of replicas for the deployment. Must be >= 0. | `true` |

```json
{
  "replicas": 20
}
```

*Response*
```json
{
  "message": "Successfully updated the replica count for the deployment"
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully updated the replica count for the deployment |
| 400 | Incorrect value for the number of replicas was set. |
| 404 | The provided namespace or deployment does not exist. |



#### **`GET`** `/health/ready`
Health check endpoint that verifies the server has access to the Kubernetes API.
Should this health check fail, the particular server should be removed from
having requests directed to it.

Note that the health check endpoints should generally not be exposed, but should
be safe in the context of internal use only.

*Parameters*
None.

*Response*
```json
{
  "status": "ready"
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | The server is able to successfully communicate with the Kubernetes API |
| 503 | The server is having issues communicating with the Kubernetes API |


#### **`GET`** `/health/live`
Health check endpoint that verifies the server is alive and responsive
(i.e. not timing out). Should this health check fail, it may signify that the
server has become unresponsive and should indicate reason to restart the server
process.

Note that the health check endpoints should generally not be exposed, but should
be safe in the context of internal use only.

*Parameters*
None.

*Response*
```json
{
  "alive": true
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | The server is alive and responsive to requests. |
| 503 | The server process may have startup or configuration issues |


### Security
#### Kubernetes RBAC Permissions
As the `sre-server` will be communicating with a Kubernetes API, it will need
sufficient RBAC priveleges to communicate with Kubernetes. At a minimum, for the
initial implementation, it will need a Kubernetes `ClusterRole` (since it needs
cluster-wide access to all namespaces to see the `Deployments`) with the
following permissions:

```yaml
rules:
  # For listing and retrieving deployments
  - apiGroups: ["apps"]
    resources:
      - deployments
    verbs:
      - get
      - list
      - watch

  # For updating the number of replicas on the deployment
  - apiGroups:
      - apps
    resources:
      - deployments/scale
    verbs:
      - patch
```

#### Authentication
We make use of mTLS for secured communications between the client and the
`sre-server`. With mTLS, both the server and client verify each other's
certificates (authenticate). For this initial design, in order to simplify
usage and development, we opted to make use of the same public and private
keys for both the server and the clients. However, for production, we will
want to use an existing internal CA (or use something like Hashicorp Vault)
to generate the TLS certifcate key pairs for each client and server for a
secure setup. For generating the TLS certificate key pairs themselves, we
recommend using Ed25519 for improved performance (and sometimes security) over
RSA and ECDSA. A key pair that can be used by both clients and the server
can be done though:

```shell
# For `-days` parameter, you can change the default to increase how long the
# certificate is valid for. We default it to a low number to catch accidental
# usage in production.
#
# For the `-addext "subjectAltName` option, we add a subject alternative name
# to prevent Go based clients from complaining about
# "certificate relies on legacy Common Name field, use SANs instead". In
# general, the industry seems to be moving towards using SANs over the common
# name.
openssl req \
  -newkey ed25519 \
  -new \
  -nodes \
  -x509 \
  -days 1 \
  -out public-certificate.pem \
  -keyout private-key.pem \
  -subj "/C=US/ST=Cali/L=Somewhere/O=Your Organization/OU=Your Unit/CN=localhost" \
  -addext "subjectAltName = DNS:localhost"
```

The `sre-server` shall enforce TLS 1.3 as the only option to enforce the
current best security practices, which is supported by the majority of
clients today. Using TLS 1.3 in Go also means that we're not able to choose or
prefer any particular cipher suite (see https://pkg.go.dev/crypto/tls#Config),
but default to the following:

-  `TLS_AES_128_GCM_SHA256`
-  `TLS_AES_256_GCM_SHA384`
-  `TLS_CHACHA20_POLY1305_SHA256`

Should we require to use TLS 1.2 to support legacy clients, then we still
recommend to use the above cipher suites as the preferred for the current
security best practices.


### Development
For any changes, make sure to add in entry into the `CHANGELOG.md`, we follow
the format as specified by
[keep a changelog](https://keepachangelog.com/en/1.1.0/). We follow
[SemVer](https://semver.org/) for versioning.


#### Tooling
After git cloning the repository, you'll need to have the following tools
installed for local development and testing:

- [go](https://go.dev/doc/install) (1.21+) - for building the server
- make - for local task automation
- [docker](https://docs.docker.com/engine/install/) - for building the server
- openssl - for generating TLS certificates for development and testing purposes (only needed if creating a new key pair, otherwise use existing set in the repository)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (or similar tool) - for running a local kubernetes cluster to test the server on

#### Running
Once the above requisite tooling has been setup, within the root of this
repository, you can run the server locally by (faster development turnaround,
but less consistency with deployment environment):

```shell
go run .
```

For running the server in a reproducible and consistent manner that is more
production like:

```shell
make local-run
```

This will build the server into a container and then run it through there.


#### Testing
For testing locally, you can run:

```shell
go test -v
```

This will run all of the go unit tests.


To run the integration tests, use:

```shell
make integration-tests
```

which will run integration tests against the local kubernetes cluster.

All of the tests will also be ran in a GitHub actions from a container for
ensuring reproducibility.


### Deliverables
#### Build
For reproducible and consistent builds, we make use of containers. You can build
it through:

```shell
make build
```
Which will create a ready to use docker container of the server.


#### Release
If you need to build a release of the server as well, you can run:

```shell
make release
```

Which will build binaries of the server in addition to the container image.

Releases shall be created through GitHub Actions on a merge to the main branch.
This requires that the `CHANGELOG.md` be appropriately updated in the pull
request that got merged so that the new version is automatically built.


#### Local Deployment
To deploy to your local Kubernetes cluster (assuming a `kind` cluster), you can
run the following command to do so:

```shell
make local-deploy
```

Which shall build the server for use with the `kind` cluster so that it is made
available for deployment via the included Helm chart.


To clean up the deployment, you can run:

```shell
make local-deploy-clean
```
