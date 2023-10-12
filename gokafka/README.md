# Go Kafka App with Observability


```bash
docker build -t gokafka:latest .
docker tag  gokafka:latest ghcr.io/kuro337/gokafka:latest
docker push ghcr.io/kuro337/gokafka:latest

```
- Go App

```bash
kubectl apply -f kafkago.yaml

kubectl delete deployment go-app-deployment

kubectl logs go-app-deployment-574b64bd-5lphc 

kubectl logs otel-collector-opentelemetry-collector-agent-pd2kz

kubectl logs otel-collector-cluster-opentelemetry-collector-545bbc4d65-cj69n
kubectl get pods

# Info for Collector Service
kubectl describe service otel-collector-cluster-opentelemetry-collector



```