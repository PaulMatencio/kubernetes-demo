kubectl delete pods healthy-monolith monolith secure-monolith
kubectl delete services monolith auth frontendmicro hello
kubectl delete deployments auth frontendmicro hello
kubectl delete secrets tls-certs
kubectl delete configmaps nginx-frontend-conf nginx-proxy-conf
kubectl delete  hpa hello

