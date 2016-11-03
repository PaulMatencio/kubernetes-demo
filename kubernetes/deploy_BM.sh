export  KUBEDEMO=$HOME/go/src/github.com/kubernetes-demo
export  kubeconf=/etc/kubernetes/admin.conf
# create secret and config
kubectl --kubeconfig $kubeconf  create secret generic tls-certs --from-file=$KUBEDEMO/kubernetes/tls/
kubectl --kubeconfig $kubeconf  create configmap nginx-frontend-conf --from-file=$KUBEDEMO/kubernetes/nginx/frontend.conf
# Deploy the microservices
kubectl  --kubeconfig $kubeconf create -f $KUBEDEMO/kubernetes/all-in-one/

