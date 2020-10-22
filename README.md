# inlets-operator

[![Build Status](https://travis-ci.com/inlets/inlets-operator.svg?branch=master)](https://travis-ci.com/inlets/inlets-operator) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/inlets/inlets-operator)](https://goreportcard.com/report/github.com/inlets/inlets-operator) [![Documentation](https://godoc.org/github.com/inlets/inlets-operator?status.svg)](http://godoc.org/github.com/inlets/inlets-operator)

Add public LoadBalancers to your local Kubernetes clusters.

In cloud-based [Kubernetes](https://kubernetes.io/) solutions, Services can be exposed as type "LoadBalancer" and your cloud provider will provision a LoadBalancer and start routing traffic, in another words: you get "Ingress" to your services from the outside world.

inlets-operator brings that same experience to your local Kubernetes cluster. The operator automates the creation of an [inlets](https://inlets.dev) exit-server on public cloud, and runs the client as a Pod inside your cluster. Your Kubernetes `Service` will be updated with the public IP of the exit-node and you can start receiving incoming traffic immediately.

## Who is this for?

This solution is for users who want to gain incoming network access (ingress) to private Kubernetes clusters. These may be running on-premises, on your laptop, within a VM or a Docker container. It even works behind NAT, and through HTTP proxies, without the need to open firewall ports. The cost of the LoadBalancer with a IaaS like DigitalOcean is around 5 USD / mo, which is several times cheaper than AWS or GCP.

## Video demo

Watch a video walk-through where we deploy an IngressController (ingress-nginx) to KinD, and then obtain LetsEncrypt certificates using cert-manager.

![Video demo](https://img.youtube.com/vi/4wFSdNW-p4Q/hqdefault.jpg)

[Try the step-by-step tutorial](https://docs.inlets.dev/#/get-started/quickstart-ingresscontroller-cert-manager?id=quick-start-expose-your-ingresscontroller-and-get-tls-from-letsencrypt-and-cert-manager)

## inlets tunnel capabilities

The operator detects Services of type LoadBalancer, and then creates a `Tunnel` Custom Resource. Its next step is to provision a small VM with a public IP on the public cloud, where it will run the inlets tunnel server. Then an inlets client is deployed as a Pod within your local cluster, which connects to the server and acts like a gateway to your chosen local service.

Pick inlets PRO or OSS.

### [inlets PRO](https://github.com/inlets/inlets-pro)

* Automatic end-to-end encryption of the control-plane using PKI and TLS
* Tunnel any TCP traffic at L4 i.e. Mongo, Postgres, MariaDB, Redis, NATS, SSH and TLS itself.
* Tunnel an IngressController including TLS termination and LetsEncrypt certs from cert-manager
* Punch out multiple ports such as 80 and 443 over the same tunnel
* Commercially licensed and supported. For cloud native operators and developers.

Heavily discounted [pricing available](https://inlets.dev/) for personal use.

### [inlets OSS](https://github.com/inlets/inlets)

* No encryption enabled for the control-plane.
* Tunnel L7 HTTP traffic.
* Punch out only one port per tunnel, port name must be: `http`
* Free, OSS, built for community developers.

If you transfer any secrets, login info, business data, or confidential information then you should use inlets PRO for its built-in encryption using TLS and PKI.

### inlets projects

Inlets is a Cloud Native Tunnel and is [listed on the Cloud Native Landscape](https://landscape.cncf.io/category=service-proxy&format=card-mode&grouping=category&sort=stars) under *Service Proxies*.

* [inlets](https://github.com/inlets/inlets) - Cloud Native Tunnel for L7 / HTTP traffic written in Go
* [inlets-pro](https://github.com/inlets/inlets-pro-pkg) - Cloud Native Tunnel for L4 TCP
* [inlets-operator](https://github.com/inlets/inlets-operator) - Public IPs for your private Kubernetes Services and CRD
* [inletsctl](https://github.com/inlets/inletsctl) - Automate the cloud for fast HTTP (L7) and TCP (L4) tunnels

## Status and backlog

Operator cloud host provisioning:

- [x] Provision VMs/exit-nodes on public cloud
  - [x] Provision to [Packet.com](https://packet.com)
  - [x] Provision to DigitalOcean
  - [x] Provision to Scaleway
  - [x] Provision to GCP
  - [x] Provision to AWS EC2
  - [x] Provision to Linode
  - [x] Provision to Azure
- [x] Provision to Civo
- [x] Publish stand-alone [Go provisioning library/SDK](https://github.com/inlets/inletsctl/tree/master/pkg/provision)

With [`inlets-pro`](https://github.com/inlets/inlets-pro) configured, you get the following additional benefits:

- [x] Automatic configuration of TLS and encryption using secured websocket `wss://` for control-port
- [x] Tunnel pure TCP traffic
- [x] Separate data-plane (ports given by Kubernetes) and control-plane (port `8132`)

Other features:

- [x] Automatically update Service type LoadBalancer with a public IP
- [x] Tunnel L7 `http` traffic
- [x] In-cluster Role, Dockerfile and YAML files
- [x] Raspberry Pi / armhf build and YAML file
- [x] ARM64 (Graviton/Odroid/Packet.com) Dockerfile/build and K8s YAML files
- [x] Control which services get a LoadBalancer using annotations
- [x] Garbage collect hosts when Service or CRD is deleted
- [x] CI with Travis and automated release artifacts
- [x] One-line installer [arkade](https://get-arkade.dev/) - `arkade install inlets-operator --help`

Backlog pending:

- [x] Feel free to request features.

## inlets-operator reference documentation for different cloud providers

Check out the reference documentation for inlets-operator to get exit-nodes provisioned on different cloud providers [here](https://docs.inlets.dev/#/tools/inlets-operator?id=inlets-operator-reference-documentation).

## Expose a service with a LoadBalancer

The LoadBalancer type is usually provided by a cloud controller, but when that is not available, then you can use the inlets-operator to get a public IP and ingress.

> The free OSS version of inlets provides a HTTP tunnel, inlets PRO can provide TCP and full functionality to an IngressController.

First create a deployment for Nginx.

For Kubernetes 1.17 and lower:

```bash
kubectl run nginx-1 --image=nginx --port=80 --restart=Always
```

For 1.18 and higher:

```bash
kubectl apply -f https://raw.githubusercontent.com/inlets/inlets-operator/master/contrib/nginx-sample-deployment.yaml
```

Now create a service of type LoadBalancer via `kubectl expose`:

```bash
kubectl expose deployment nginx-1 --port=80 --type=LoadBalancer
kubectl get svc

kubectl get tunnel/nginx-1-tunnel -o yaml

kubectl logs deploy/nginx-1-tunnel-client
```

Check the IP of the LoadBalancer and then access it via the Internet.

## Notes on OSS inlets

inlets PRO can tunnel multiple ports, but inlets OSS is set to take the first port named "http" for your service. With the OSS version of inlets (see example with OpenFaaS), make sure you give the `port` a `name` of `http`, otherwise a default of `80` will be used incorrectly.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: gateway
  namespace: openfaas
  labels:
    app: gateway
spec:
  ports:
    - name: http
      port: 8080
      protocol: TCP
      targetPort: 8080
      nodePort: 31112
  selector:
    app: gateway
  type: LoadBalancer
```

## Annotations and ignoring services

By default the operator will create a tunnel for every LoadBalancer service.

There are two ways to override the behaviour:

1) Create LoadBalancers for every service, unless annotated

  To ignore a service such as `traefik` type in: `kubectl annotate svc/traefik -n kube-system dev.inlets.manage=false`

2) Create LoadBalancers only annotated services

  You can also set the operator to ignore the services by default and only manage them when the annotation is true. `dev.inlets.manage=true` with the flag `-annotated-only`

## Monitor/view logs

The operator deployment is in the `kube-system` namespace.

```sh
kubectl logs deploy/inlets-operator -n kube-system -f
```

## Running on a Raspberry Pi

Use the same commands as described in the section above.

> There used to be separate deployment files in `artifacts` folder called `operator-amd64.yaml` and `operator-armhf.yaml`.
> Since version `0.2.7` Docker images get built for multiple architectures with the same tag which means that there is now just one deployment file called `operator.yaml` that can be used on all supported architecures.

# Provider Pricing

The host [provisioning code](https://github.com/inlets/inletsctl/tree/master/pkg/provision) used by the inlets-operator is shared with [inletsctl](https://github.com/inlets/inletsctl), both tools use the configuration in the grid below.

These costs need to be treated as an estimate and will depend on your bandwidth usage and how many hosts you decide to create. You can at all times check your cloud provider's dashboard, API, or CLI to view your exit-nodes. The hosts provided have been chosen because they are the absolute lowest-cost option that the maintainers could find.

| Provider                                                           | Price per month | Price per hour |     OS image | CPU | Memory | Boot time |
| ------------------------------------------------------------------ | --------------: | -------------: | -----------: | --: | -----: | --------: |
| [Google Compute Engine](https://cloud.google.com/compute)                                          |         *  ~\$4.28 |       ~\$0.006 | Debian GNU Linux 9 (stretch) | 1 | 614MB | ~3-15s |
| [Packet](https://www.packet.com/cloud/servers/t1-small/)           |           ~\$51 |         \$0.07 | Ubuntu 16.04 |   4 |    8GB | ~45-60s  |
| [Digital Ocean](https://www.digitalocean.com/pricing/#Compute)     |             \$5 |      ~\$0.0068 | Ubuntu 16.04 |   1 |  512MB | ~20-30s  |
| [Scaleway](https://www.scaleway.com/en/pricing/#virtual-instances) |           2.99€ |         0.006€ | Ubuntu 18.04 |   2 |    2GB | 3-5m      |

* The first f1-micro instance in a GCP Project (the default instance type for inlets-operator) is free for 720hrs(30 days) a month 

## Contributing

Contributions are welcome, see the [CONTRIBUTING.md](CONTRIBUTING.md) guide.

## Similar projects / products and alternatives

- [inlets pro](https://github.com/inlets/inlets-pro) - L4 TCP tunnel, which can tunnel any TCP traffic with automatic, built-in encryption. Kubernetes-ready with Docker images and YAML manifests. 
- [inlets](https://inlets.dev) - inlets provides an L7 HTTP tunnel for applications through the use of an exit node, it is used by the inlets operator. Encryption can be configured separately.
- [metallb](https://github.com/danderson/metallb) - open source LoadBalancer for private Kubernetes clusters, no tunnelling.
- [Cloudflare Argo](https://www.cloudflare.com/en-gb/products/argo-tunnel/) - paid SaaS product from Cloudflare for Cloudflare customers and domains - K8s integration available through Ingress
- [ngrok](https://ngrok.com) - a popular tunnelling tool, restarts every 7 hours, limits connections per minute, paid SaaS product with no K8s integration available

## Author / vendor

inlets and the inlets-operator are brought to you by [OpenFaaS Ltd](https://www.openfaas.com) and [Alex Ellis](https://www.alexellis.io/).
