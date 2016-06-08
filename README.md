#  Scalable microservices with Kubernetes and containers
#  Built on top of the udacity  ud615 course)

You will be doing

* Provision a complete Kubernetes using [Google Container Engine](https://cloud.google.com/container-engine)
* Deploy and manage  containers (Docker) using kubectl

Kubernetes Version: 1.2.2

## Description

Kubernetes is all about applications and  you will be utilizing the Kubernetes API to deploy, manage, and upgrade applications.

App is an example 12 Facter application that we'll be using. You will be working with the following Docker images:

* [pmatencio60/monolith](https://hub.docker.com/r/pmatencio60/monolith) - Monolith includes auth and hello services.
* [pmatencio60/auth](https://hub.docker.com/r/pmatencio60/auth) - Auth microservice. Generates JWT tokens for authenticated users.
* [pmatencio60/hello](https://hub.docker.com/r/pmatencio60/hello) - Hello microservice. Greets authenticated users.
* [ngnix](https://hub.docker.com/_/nginx) - Frontend to the auth and hello services.

## Links

  * [Kubernetes](http://googlecloudplatform.github.io/kubernetes)
  * [gcloud Tool Guide](https://cloud.google.com/sdk/gcloud)
  * [Docker](https://docs.docker.com)
  * [etcd](https://coreos.com/docs/distributed-configuration/getting-started-with-etcd)
  * [nginx](http://nginx.org)


## Run the Demo from code to deployment

- Prerequisites
- Building and tag container image (Docker)
- Pushing tag image to containers Repository (Docker Hub)
- Creating kubernetes cluster ( 3 nodes in  GCE)
- Create a load balancer to the cluster
- Creating  secret (https micro services)
- Creating configmap (decouple configuration artifacts from the image content)
- Creating deployments
- Creating services
- Scaling as needed( replicasets)/ autoscaling
- Rolling update
- Rolling back

## Prerequisites
* Get a docker hub account
* Get  a Google cloud platform account
* Enable the Google Compute Engine and Google Container engine API

# Start your Google Cloud shell
- Create a new Google Cloud platform project
- Clone the Kubernetes-demo github  repository in your cloud shell
  - cd ~/go/src; mkdir github; cd github
  - git pull  https://github.com/PaulMatencio/kubernetes-demo
  - cd kubernetes-demo/app

# Building docker images

- Download and install  the latest go release  - https://golang.org/dl/-
  - export GOPATH=~/go;
  - export GOROOT=<go installation path>  - /usr/local/go -
- cd  ~/go/src/github/kubernetes-demo/app/monolith
  - go build --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"'
  - sudo docker build -t monolith:1.0.0   .
- cd  ../auth
  - go build --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"'
  - sudo docker build -t auth:1.0.0  .
- cd ../hello
  - go build --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"'
  - sudo docker build -t  hello:1.0.0  

# Tagging docker images

- List the docker images :  sudo docker images
- Tag the  auth:1.0.0  image
  - cd ~/go/src/github.com/kubernetes-demo/app/auth
  - sudo docker tag auth:1.0.0  <your-hub-username>/auth:1.0.0
- Tag the hello:1.0.0 image
  - cd ../hello
  - sudo docker tag  hello:1.0.0  <your-hub-username>/hello:1.0.0
- Tag the monolith image
  - cd ../monolith
  - sudo docker tag  monolith:1.0.0  <your-hub-username>/monolith:1.0.0

# Pushing docker images to the docker hub
- Login to the docker hub
  - docker login --username=your-hub-username --email=your-email@company.com

- Push tagged docker images to docker hub
  - docker push < your-hub-username>/auth:1.0.0
  - docker push  <your-hub-username>/hello:1.0.0
  - docker push <your-hub-username>/monolith:1.0.0
