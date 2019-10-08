# inlets-operator

Get a Kubernetes LoadBalancer where you never thought it was possible.

In cloud-based [Kubernetes](https://kubernetes.io/) solutions, Services can be exposed as type "LoadBalancer" and your cloud provider will provision a LoadBalancer and start routing traffic, in another word: you get ingress to your service.

inlets-operator brings that same experience to your local Kubernetes or k3s cluster (k3s/k3d/minikube/microk8s/Docker Desktop/KinD). The operator automates the creation of an [inlets](https://inlets.dev) exit-node on public cloud, and runs the client as a Pod inside your cluster. Your Kubernetes `Service` will be updated with the public IP of the exit-node and you can start receiving incoming traffic immediately.

## Who is this for?

This solution is for users who want to gain incoming network access (ingress) to their private Kubernetes clusters running on their laptops, VMs, within a Docker container, on-premises, or behind NAT. The cost of the LoadBalancer with a IaaS like DigitalOcean is around 5 USD / mo, which is 10 USD cheaper than an AWS ELB or GCP LoadBalancer.

Whilst 5 USD is cheaper than a "Cloud Load Balancer", this tool is for users who cannot get incoming connections due to their network configuration, not for saving money vs. public cloud.

## Status and backlog

This version of the inlets-operator is a early proof-of-concept, but it builds upon inlets, which is stable and widely used.

Backlog completed:
- [x] Provision VMs/exit-nodes on public cloud
- [x] Provision to [Packet.com](https://packet.com)
- [x] Provision to DigitalOcean
- [x] Automatically update Service type LoadBalancer with a public IP
- [x] Tunnel L7 `http` traffic
- [x] In-cluster Role, Dockerfile and YAML files
- [x] Raspberry Pi / armhf build and YAML file
- [ ] Ignore Services with `dev.inlets.manage: false` annotation

Backlog pending:
- [ ] Garbage collect hosts when CRD is deleted
- [ ] CI with Travis (use openfaas-incubator/openfaas-operator as a sample)
- [ ] ARM64 (Graviton/Odroid/Packet.com) Dockerfile/build and K8s YAML files
- [ ] Automate `wss://` for control-port
- [ ] Move control-port and `/tunnel` endpoint to high port i.e. `31111`
- [ ] Provision to EC2
- [ ] Provision to GCP
- [ ] Tunnel any `tcp` traffic (using `inlets-pro`)

Inlets tunnels HTTP traffic at L7, so the inlets-operator can be used to tunnel HTTP traffic. A new project I'm working on called inlets-pro tunnels any TCP traffic at L4 i.e. Mongo, Redis, NATS, SSH, TLS, whatever you like.

## Author

inlets and inlets-operator are brought to you by [Alex Ellis](https://twitter.com/alexellisuk). Alex is a [CNCF Ambassador](https://www.cncf.io/people/ambassadors/) and the founder of [OpenFaaS](https://github.com/openfaas/faas/).

If you like this project, then join dozens of other developers by Sponsoring Alex and his OSS work through [GitHub Sponsors](https://github.com/users/alexellis/sponsorship) today.

## Video demo

This video demo shows a single-node VM running on k3s on Packet.com, and the inlets exit node also being provisioned on Packet's infrastructure.

[![https://img.youtube.com/vi/LeKMSG7QFSk/0.jpg](https://img.youtube.com/vi/LeKMSG7QFSk/0.jpg)](https://www.youtube.com/watch?v=LeKMSG7QFSk&amp=&feature=youtu.be)

See an alternative video showing my cluster running with KinD on my Mac and the exit node being provisioned on DigitalOcean:

* [KinD & DigitalOcean](https://youtu.be/c6DTrNk9zRk).

## Step-by-step tutorial

[Try the step-by-step tutorial](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)

## Running in-cluster, using DigitalOcean for the exit node

You can also run the operator in-cluster, a ClusterRole is used since Services can be created in any namespace, and may need a tunnel.

```sh
# Create a secret to store the access token

kubectl create secret generic inlets-access-key \
  --from-literal inlets-access-key="$(cat ~/Downloads/do-access-token)"

# Apply the operator deployment and RBAC role
kubectl apply -f ./artifacts/operator-rbac.yaml
kubectl apply -f ./artifacts/operator-amd64.yaml
```

## Running on a Raspberry Pi (armhf), using DigitalOcean for the exit node

To get a LoadBalancer for services running on your Raspberry Pi, use the armhf deployment file:

```sh
# Create a secret to store the access token

kubectl create secret generic inlets-access-key \
  --from-literal inlets-access-key="$(cat ~/Downloads/do-access-token)"

# Apply the operator deployment and RBAC role
kubectl apply -f ./artifacts/operator-rbac.yaml
kubectl apply -f ./artifacts/operator-armhf.yaml
```

## Run the Go binary with Packet.com

Assuming you're running a local cluster with [KinD](https://github.com/kubernetes-sigs/kind):

Sign up to [Packet.com](https://packet.com) and get an access key, save it in `~/packet-token`

```sh
kubectl apply ./aritifacts/crd.yaml

export PACKET_PROJECT_ID=""	# Populate from dashboard

export GOPATH=$HOME/go/
go get -u github.com/alexellis/inlets-operator
cd $GOPATH/github.com/alexellis/inlets-operator

go get

go build && ./inlets-operator  --kubeconfig "$(kind get kubeconfig-path --name="kind")" --access-key=$(cat ~/packet-token) --project-id="${PACKET_PROJECT_ID}"
```

## Run the Go binary with DigitalOcean

Assuming you're running a local cluster with [KinD](https://github.com/kubernetes-sigs/kind):

Sign up to [DigitalOcean.com](https://DigitalOcean.com) and get an access key, save it in `~/do-access-token`.

```sh
kubectl apply ./aritifacts/crd.yaml

export GOPATH=$HOME/go/
go get -u github.com/alexellis/inlets-operator
cd $GOPATH/github.com/alexellis/inlets-operator

go get

go build && ./inlets-operator  --kubeconfig "$(kind get kubeconfig-path --name="kind")" --access-key=$(cat ~/do-access-token) --provider digitalocean
```

# Monitor/view logs

```sh
kubectl logs deploy/inlets-operator -f
```

## Get a LoadBalancer provided by inlets

```sh
kubectl run nginx-1 --image=nginx --port=80 --restart=Always
kubectl run nginx-2 --image=nginx --port=80 --restart=Always

kubectl expose deployment nginx-1 --port=80 --type=LoadBalancer
kubectl expose deployment nginx-2 --port=80 --type=LoadBalancer

kubectl get svc

kubectl get tunnel nginx-tunnel-1 -o yaml

kubectl get svc

kubectl logs deploy/nginx-1-tunnel-client
```

Check the IP of the LoadBalancer and then access it via the Internet.

Example with OpenFaaS, make sure you give the `port` a `name` of `http`, otherwise a default of `80` will be used incorrectly.

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

To ignore a service such as `traefik` type in: `kubectl annotate svc/traefik -n kube-system dev.inlets.manage=false`

# Provider Pricing

| Provider                                                       | Price per month | Price per hour |  OS image    | CPU | Memory |
| -------------------------------------------------------------- | --------------: | -------------: | -----------: | --: | -----: |
| [Packet](https://www.packet.com/cloud/servers/t1-small/)       | ~$51            | $0.07          | Ubuntu 16.04 | 4   | 8GB    |
| [Digital Ocean](https://www.digitalocean.com/pricing/#Compute) | $5              | ~$0.0068       | Ubuntu 16.04 | 1   | 512MB  |

## Contributing

Contributions are welcome, see the [CONTRIBUTING.md](CONTRIBUTING.md) guide.

## Similar projects / products and alternatives

* [metallb](https://github.com/danderson/metallb) - open source LoadBalancer for private Kubernetes clusters, no tunnelling.
* [inlets](https://inlets.dev) - inlets provides an L7 HTTP tunnel for applications through the use of an exit node, it is used by the inlets operator
* inlets pro - L4 TCP tunnel, which can tunnel any TCP traffic and is on the roadmap for the inlets-operator
* [Cloudflare Argo](https://www.cloudflare.com/en-gb/products/argo-tunnel/) - paid SaaS product from Cloudflare for Cloudflare customers and domains - K8s integration available through Ingress
* [ngrok](https://ngrok.com) - a popular tunnelling tool, restarts every 7 hours, limits connections per minute, paid SaaS product with no K8s integration available

