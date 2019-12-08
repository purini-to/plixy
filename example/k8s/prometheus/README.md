```bash
kubectl create namespace prometheus
helm install -f ./config.yaml prometheus stable/prometheus --namespace prometheus
helm install -f ./adapter-config.yaml prometheus-adapter stable/prometheus-adapter --namespace prometheus
```
