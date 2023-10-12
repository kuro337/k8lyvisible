# Collector

- We have apps - and using OpenTelemetry SDK - they can send logs directly to a Backend , or stdout , etc.

- We can have a Collector in between - so apps send data to Collector - then Collector sends to the backend.

- We can have a Single Central Collector - or collectors Running as a SideCar, etc.

- Deploying the Collector 

- https://opentelemetry.io/docs/kubernetes/getting-started/

```bash
# Add helm chart so we can add it later
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts

# Deploy DaemonSet for Collector
# Specify values.yaml for Exporter settings , etc.
helm install otel-collector open-telemetry/opentelemetry-collector --values collector-daemonset.yaml
helm uninstall otel-collector 


helm install otel-collector open-telemetry/opentelemetry-collector --values collector-daemon-loki.yaml



```

- Deployment Collector

```bash
# Deploy Collector 2nd component
helm install otel-collector-cluster open-telemetry/opentelemetry-collector --values collector-deployment.yaml

helm uninstall otel-collector-cluster

# Getting endpoints
kubectl get service
kubectl describe service otel-collector-cluster-opentelemetry-collector

# Getting applied config
kubectl get daemonset otel-collector-opentelemetry-collector-agent  -o yaml

# Config Map to presets - to show receivers 
kubectl get configmap otel-collector-opentelemetry-collector-agent -o yaml


```

- Testing endpoints like gRPC connections

```bash
# Test grpc connection
kubectl run grpcurl-test --image=fullstorydev/grpcurl --restart=Never --command -- /bin/grpcurl -plaintext otel-collector-cluster-opentelemetry-collector.default.svc.cluster.local:4317 list

kubectl get pods
kubectl logs grpcurl-test   
kubectl delete pod grpcurl-test 

```

- For the app

```bash
We have Collector that uses FileLog to get logs from stdout

For Traces and Metrics - send them to the Central Collector

gradle assemble
env \
OTEL_SERVICE_NAME=dice-server \
OTEL_TRACES_EXPORTER=otlp \
OTEL_EXPORTER_OTLP_ENDPOINT=http://<otel-collector-service-name>:4317 \
OTEL_METRICS_EXPORTER=otlp \
OTEL_EXPORTER_OTLP_METRIC_ENDPOINT=http://<otel-collector-service-name>:4317 \
java -jar ./build/libs/java-simple.jar


```

- Debug

```bash
# Checking Daemonset logs for Collector
kubectl logs otel-collector-opentelemetry-collector-agent-t2d5k 

kubectl logs otel-collector-cluster-opentelemetry-collector-545bbc4d65-cj69n

2023-10-11T22:38:47.041Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 47, "metrics": 519, "data points": 565}
2023-10-11T22:38:51.657Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 1, "metrics": 38, "data points": 48}
2023-10-11T22:38:57.703Z        info    fileconsumer/file.go:194        Started watching file   {"kind": "receiver", "name": "filelog", "data_type": "logs", "component": "fileconsumer", "path": "/var/log/pods/default_go-app-deployment-574b64bd-9vhph_454a9f0d-eaf6-4841-b2fa-f05deb073034/go-app/0.log"}

kubectl logs otel-collector-cluster-opentelemetry-collector-545bbc4d65-cj69n

kubectl logs loki-0

```

- Deploying an App to test Tracing

```bash
# Get endpoint to specify for our App's Env Variable

kubectl describe pod otel-collector-opentelemetry-collector-agent-wsmm2
kubectl get services

# We can see service is otel-collector-cluster.... - specify this in the java-app.yaml

# Setting this env variable the OTEL SDK knows where to send Traces

kubectl apply -f java-app.yaml

kubectl delete deployment dice-server  

kubectl logs dice-server-7846894584-x67bk 


kubectl exec -it dice-server -- /bin/sh


# Expose App - then run curl from the debug container

kubectl expose deployment dice-server --type=ClusterIP --port=8080
kubectl get service dice-server


# Debug Container
kubectl run curl-debug --image=curlimages/curl:latest --rm -it --restart=Never -- /bin/sh
curl http://dice-server:8080/rolldice?rolls=12

# Check Logs
kubectl logs dice-server-558dc55496-tn9kr 

```

- Example yaml for App

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dice-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dice-server
  template:
    metadata:
      labels:
        app: dice-server
    spec:
      containers:
        - name: dice-server
          image: ghcr.io/kuro337/jar-trace-simple-tempo:latest
          env:
            - name: OTEL_SERVICE_NAME
              value: "dice-server"
            - name: OTEL_EXPORTER_OTLP_ENDPOINT
              value: "http://otel-collector-cluster-opentelemetry-collector:4317"
            - name: OTEL_EXPORTER_OTLP_PROTOCOL
              value: "http/protobuf"
            - name: OTEL_EXPORTER_OTLP_INSECURE
              value: "true"

          ports:
            - containerPort: 8080
      imagePullSecrets:
        - name: ghcr-secret

```