# Kube No Trouble - `kubent`

__*Easily check your cluster for use of deprecated APIs*__

Kubernetes 1.16 is slowly starting to roll out, not only across various managed
Kubernetes offerings, and with that come a lot of API deprecations[1][1].

*Kube No Trouble (__`kubent`__)* is a simple tool to check whether you're using any
of these API versions in your cluster and therefore should upgrade your
workloads first, before upgrading your Kubernetes cluster.

This tool will be able to detect deprecated APIs depending on how you deploy
your resources, as we need the original manifest to be stored somewhere. In
particular following tools are supported:
- **file**    - local manifests in YAML or JSON
- **kubectl** - uses the `kubectl.kubernetes.io/last-applied-configuration` annotation
- **Helm v2** - uses Tiller manifests stored in K8s Secrets or ConfigMaps
- **Helm v3** - uses Helm manifests stored as Secrets or ConfigMaps directly in individual namespaces

[1]: https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/

**Additional resources:**
- Blog post on K8s deprecated APIs and introduction of kubent: [Kubernetes: How to automatically detect and deal with deprecated APIs][2]

[2]: https://blog.doit-intl.com/kubernetes-how-to-automatically-detect-and-deal-with-deprecated-apis-f9a8fc23444c

## Install

Run `sh -c "$(curl -sSL https://git.io/install-kubent)"`.

*(The script will download latest version and unpack to `/usr/local/bin`).*

Or download the
[latest release](https://github.com/doitintl/kube-no-trouble/releases/latest)
for your platform and unpack manually.


## Usage

Configure Kubectl's current context to point to your cluster, `kubent` will
look for the kube `.config` file in standard locations (you can point it to custom
location using the `-k` switch). 

**`kubent`** will collect resources from your cluster and report on found issuses.

*Please note that you need to have sufficient permissions to read Secrets in the
cluster in order to use `Helm*` collectors.*

```sh
$./kubent
6:25PM INF >>> Kube No Trouble `kubent` <<<
6:25PM INF Initializing collectors and retrieving data
6:25PM INF Retrieved 103 resources from collector name=Cluster
6:25PM INF Retrieved 132 resources from collector name="Helm v2"
6:25PM INF Retrieved 0 resources from collector name="Helm v3"
6:25PM INF Loaded ruleset name=deprecated-1-16.rego
6:25PM INF Loaded ruleset name=deprecated-1-20.rego
__________________________________________________________________________________________
>>> 1.16 Deprecated APIs <<<
------------------------------------------------------------------------------------------
KIND         NAMESPACE     NAME                    API_VERSION
Deployment   default       nginx-deployment-old    apps/v1beta1
Deployment   kube-system   event-exporter-v0.2.5   apps/v1beta1
Deployment   kube-system   k8s-snapshots           extensions/v1beta1
Deployment   kube-system   kube-dns                extensions/v1beta1
__________________________________________________________________________________________
>>> 1.20 Deprecated APIs <<<
------------------------------------------------------------------------------------------
KIND      NAMESPACE   NAME           API_VERSION
Ingress   default     test-ingress   extensions/v1beta1
```

### Arguments

You can list all the configuration options available using `--help` switch:
```sh
$./kubent -h
Usage of ./kubent:
  -c, --cluster             enable Cluster collector (default true)
  -d, --debug               enable debug logging
  -f, --filename strings    manifests to check
      --helm2               enable Helm v2 collector (default true)
      --helm3               enable Helm v3 collector (default true)
  -k, --kubeconfig string   path to the kubeconfig file (default "/Users/stepan/.kube/config")
  -o, --output string       output format - [text|json] (default "text")
```

### Use in CI

`kubent` will return `0` exit code if the program succeeds, even if it finds
deprecated resources, and non-zero exit code if there is an error during
runtime. Because all info output goes to stderr, it's easy to check in shell if
any issues were found:

```shell
test -z "$(kubent)"                 # if stdout output is empty, means no issuse were found
                                    # equivalent to [ -z "$(kubent)" ]
```

It's actually better so split this into two steps, in order to differentiate
between runtime error and found issues:

```shell
if ! OUTPUT="$(kubent)"; then       # check for non-zero return code first
  echo "kubent failed to run!"
elif [ -n "${OUTPUT}" ]; then       # check for empty stdout 
  echo "Deprecated resources found"
fi
```

Alternatively, use the json output and smth. like `jq` to check if the result is
empty:

```
kubent -o json | jq -e 'length == 0'
```

## Development

The simplest way to build `kubent` is:

```sh
# Clone the repository
git clone https://github.com/doitintl/kube-no-trouble.git
cd kube-no-trouble/
# We require statik for generating static embedded files
go get github.com/rakyll/statik
# Generate
go generate
# Build
go build -o bin/kubent cmd/kubent/main.go
```

Otherwise there's `Makefile`
```sh
$ make
make
all                            Cean, build and pack
help                           Prints list of tasks
build                          Build binary
generate                       Go generate
pack                           Pack binaries with upx
release-artifacts              Create release artifacts
clean                          Clean build artifacts
```

## Issues and Contributions

Please open any issues and/or PRs against github.com/doitintl/kube-no-trouble repository.

Feedback and contributions are always welcome!

### Todo:

Some future features ideas:
- Advice on correct replacement API version ?
- Output - pdf/html?
- Tests
- Lint
