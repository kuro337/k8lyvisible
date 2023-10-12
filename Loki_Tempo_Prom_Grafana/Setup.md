# Loki + Grafana

- Setting up a cluster with Grafana (viz) + Loki (logs) +  + Tempo (traces)

```bash
k3d create cluster 

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

# Setup Tempo, Loki , Grafana

helm upgrade --install tempo grafana/tempo-distributed -f tempo.yaml

helm upgrade --install  loki grafana/loki -f loki.yaml

helm upgrade --install prometheus prometheus-community/prometheus -f prometheus.yaml

helm upgrade --install grafana grafana/grafana -f grafana.yaml

# Save output for Grafana
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=grafana" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace default port-forward $POD_NAME 3000

# Testing Loki 
kubectl run curl-temp-pod --image=curlimages/curl --restart=Never --command -- sleep infinity

kubectl exec -it curl-temp-pod -- sh

curl localhost:8080

kubectl exec curl-temp-pod -- curl -G -H "X-Scope-OrgID: 1" http://loki:3100/loki/api/v1/query --data-urlencode 'query={label="testLabel"}'


kubectl exec -i curl-temp-pod -- curl -XPOST -H "Content-Type: application/json" -H "X-Scope-OrgID: 1" --data \
'{"streams": [{"stream": {"label": "testLabel"},"values": [["'$(date +%s%N)'", "This is a test log message"]] }] }' \
http://loki:3100/loki/api/v1/push

```

- Querying Logs in Grafana Loki

```sql

{exporter="OTLP"} |~ "go-app"

{exporter="OTLP"} |~ "k8s.deployment.name\":\"go-app-deployment\""

```

-- Refer to Collector/ for setting up the Collector to scrape logs from Deployments and capture Cluster Metadata too.

