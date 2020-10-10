# kube-ns

Configure Kubernetes resources, such as network policies, on creating or updating namespaces through a custom Kubernetes controller

## Usage

TBC

## Install

### Local Installation

#### Using go package installer:

```console
# Download and install kubens
$ go get -u github.com/camilocot/kube-ns

# Enable default network policy creation
$ kube-ns config add netpol --enabled

# start kubens server
$ kube-ns

INFO[0000] Starting kubens controller
INFO[0000] Processing namespace default added
INFO[0000] Processing namespace kube-public added
INFO[0000] Processing namespace kube-node-lease added
INFO[0000] Processing namespace kube-system added
INFO[0000] kubens controller synced and ready

```

## Configure

TBC

## Build

### Using go

Clone the repository anywhere:

```bash
git clone https://github.com/camilocot/kube-ns.git
cd kube-ns
go build
```

or

You can also use the Makefile directly:

```bash
make build
```

## Todo

- [ ] Quota configuration

## ACKNOWLEDGEMENTS

Most of the code from this project comes from [kube-watch](https://github.com/bitnami-labs/kubewatch) project
