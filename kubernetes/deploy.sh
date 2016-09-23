export  KUBEDEMO=$HOME/go/src/github.com/kubernetes-demo
export  GC=$HOME/google-cloud-sdk/
# create secret and config
${GC}/bin/kubectl create secret generic tls-certs --from-file=$KUBEDEMO/kubernetes/tls/
${GC}/bin/kubectl create configmap nginx-frontend-conf --from-file=$KUBEDEMO/kubernetes/nginx/frontend.conf
# Deploy the microservices
$GC/bin/kubectl  create -f $KUBEDEMO/kubernetes/all-in-one/

