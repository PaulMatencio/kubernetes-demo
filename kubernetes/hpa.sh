# autoscale horizontal pods
if [ $# -gt 0 ]
then 
  kubectl --kubeconfig $1 autoscale deployment hello  --cpu-percent=50 --min=2 --max=5
  kubectl --kubeconfig $1 get hpa
else 
  kubectl autoscale deployment hello  --cpu-percent=50 --min=2 --max=5
  kubectl get hpa
fi
