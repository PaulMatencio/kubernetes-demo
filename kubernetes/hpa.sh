# autoscale horizontal pods
kubectl autoscale deployment hello  --cpu-percent=5 --min=2 --max=5
kubectl get hpa
