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
- **kubectl** - uses the `kubectl.kubernetes.io/last-applied-configuration` annotation
- **Helm v2** - uses Tiller manifests stored in K8s Secrets or ConfigMaps
- **Helm v3** - uses Helm manifests stored as Secrets or ConfigMaps directly in individual namespaces

[1]: https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/

## Install

Download [latest
release](https://github.com/doitintl/kube-no-trouble/releases/latest) for your
platform, and unpack - `tar -xvzf kubent-*.tar.gz`.

## Usage

Configure Kubectl's current context to point to your cluster, `kubent` will
look for the kube `.config` file in standard locations (you can point it to custom
location using the `-k` switch). 

**`kubent`** will collect resources from your cluster and report on found issuses:

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
      --helm2               enable Helm v2 collector (default true)
      --helm3               enable Helm v3 collector (default true)
  -k, --kubeconfig string   path to the kubeconfig file (default "/Users/stepan/.kube/config")
  -o, --output string       output format - [text|json] (default "text")
```

## Issues and Contributions

Please open any issues and/or PRs against github.com/doitintl/kube-no-trouble repository.

Feedback and contributions are always welcome!

### Todo:

Some future features ideas:
- Input from files
- Advice on correct replacement API version ?
- Output - pdf/html?
- Tests
- Lint

