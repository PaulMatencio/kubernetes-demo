#  Scalable microservices with Kubernetes and containers
###  Built on top of the udacity  ud615 course

You will be doing ( check the run the demo below)

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


# Run the demo from code to deployment

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

## Start your Google Cloud shell

- Create a new Google Cloud platform project
- Clone the Kubernetes-demo github  repository in your cloud shell
  - cd ~/go/src; mkdir github; cd github
  - git pull  https://github.com/PaulMatencio/kubernetes-demo
  - cd kubernetes-demo/app

## Building docker images

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

## Tagging docker images

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

## Pushing docker images to the docker hub

- Login to the docker hub
  - docker login --username=your-hub-username --email=your-email@company.com

- Push tagged docker images to docker hub
  - docker push < your-hub-username>/auth:1.0.0
  - docker push  <your-hub-username>/hello:1.0.0
  - docker push <your-hub-username>/monolith:1.0.0

## Creating the dem0-k0 container cluster on GCE

- List the zones:  gcloud compute zones list
- Create a container cluster demo-k0 with 3 nodes ( default template)
  - gcloud container clusters --zone  europe-west1-b   create  demo-k0
- List the nodes : gcloud compute instances list
- Delete the container cluster demo-k0 in zone europe-west1-b
  - gcloud container clusters --zone  europe-west1-b delete  demo-k0
- Resize a container cluster
  - gcloud container clusters resize europe-west1-b  --size 4
  - gcloud container clusters resize europe-west1-b  --node-pool xx --size 4

## Exposing nginx pod as being load balancer service to the cluster  

- Run nginx in a pod
  - kubectl run nginx --image=nginx:1.10.0
- Get the list of running pods:  kubectl get pods
- Expose the running  nginx as being  a load balancer service to the cluster
  - kubectl expose deployment nginx --port 80 --type LoadBalancer
- Get the list of running services
  - kubectl get services

## Creating  secret and configmap objects

- cd ~/go/src/github/kubernetes-demo/kubernetes

- Create secret objects  - Objects of type secret are intended to hold sensitive information, such as passwords, OAuth tokens, and ssh keys -
  - kubectl create secret generic tls-certs --from-file=tls/   
  - kubectl get secret
  - kubectl describe secret tls-certs

- Create the configmap - configuration artifacts that are decoupled from image content in order to keep containerized applications portable
  - kubectl create configmap  nginx-frontend-conf  --from-file=nginx/frontend.conf
  - kubectl describe configmap nginx-frontend-conf

## Deploying auth and hello microservices

- Replace  the auth and hello deployments with your own docker images
  - vim  deployments/auth.yaml
  - vim deployments/hello.yaml

- Create the auth and hello deployments
  - kubectl  create -f deployments/auth.yaml;kubectl create -f deployments/hello.yaml

- Create the auth and hello services
  - kubectl  create -f  services/auth.yaml; kubectl  create -f  services/hello.yaml

- Check the deployments, pods, replicasets and services
  - kubectl get deployments
  - kubectl get pods
  - kubectl get services
  - kubectl get replicasets

## Deploy the frontend application

- Create the frontend deployment and service
  - kubectl create -f deployments/frontend.yaml
  - kubectl  create -f  services/frontend.yaml
- Get the external IP address of the frontend service
  - kubectl get services frontend

- Read the container log
  - Kubectl logs <pod> <container-name>

## Using the microservices

- curl -k https://<frontend external IP address>
    {"message":"Hello"}

- curl -k https://<frontend external IP address>/secure
    authorization failed, you must login first

- curl -k https://<frontend external IP address>/login - u user   (password: password)

- TOKEN=$(curl -k https://<frontend external IP address>/login -u user |   jq -r “.token”)

- curl -k  -H “Authorization:Bearer $TOKEN”  http:// https://<frontend external IP address>/secure

## Scale and Roll back
- Scaling out  the microservices
  - kubectl scale deployment/frontend   --replicas=2
  - kubectl scale deployment/auth  --replicas=2
  - kubectl scale deployment/hello   --replicas=5

- Scaling down
  - kubectl scale deployment/hello   --replicas=3

- Rollback
  - kubectl edit deployment/auth( change the docker image)
  - kubectl  apply    -f deployments/auth.yaml
  - kubectl rollout undo  deployment/auth   

## Horizontal pods autoscalling

Horizontal pod autoscaling allows to automatically scale the number of pods in a replication controller, deployment or replica set based on observed CPU utilization. In the future also other metrics will be supported.   

- kubectl create  -f  horizontalPodAutoscaler/hello.yaml
- kubectl delete  horizontalpodautoscalers  hello
- kubectl  describe  horizontalpodautoscalers hello

- start another terminal ( Google console)
- Run /bin/sh in a pod and call the hello microservice in a loop
  - kubectl run -i --tty load-generator --image=busybox /bin/sh
  - while true; do wget -q -O- http://hello.default.svc.cluster.local; done
  - <ctrl> + C to stop  the load
  - Resume : kubectl attach load-generator-xxxxx  -c load-generator -i -t

- kubectl  delete deployments load-generator  

## Drain/cordon/uncordon  a kubernetes  node
- gcloud  compute instances list

- kubectl drain/cordon  <gke-instance>
  - kubectl get pods
  - kubectl get services
  - kubectl get replicasets

- kubectl uncordon <gke-instance>
  - kubectl get pods
  - kubectl get services
  - kubectl get replicasets
