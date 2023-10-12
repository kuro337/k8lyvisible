```bash
██╗  ██╗ █████╗                                                                   
██║ ██╔╝██╔══██╗                                                                  
█████╔╝ ╚█████╔╝                                                                  
██╔═██╗ ██╔══██╗                                                                  
██║  ██╗╚█████╔╝██╗                                                               
██╗  ╚═╝ ╚════╝ ╚═╝                                                               
╚██████╗ ██████╗ ███████╗███████╗██████╗ ██╗   ██╗ █████╗ ██████╗ ██╗     ███████╗
██╔═══██╗██╔══██╗██╔════╝██╔════╝██╔══██╗██║   ██║██╔══██╗██╔══██╗██║     ██╔════╝
██║   ██║██████╔╝███████╗█████╗  ██████╔╝██║   ██║███████║██████╔╝██║     █████╗  
██║   ██║██╔══██╗╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══██║██╔══██╗██║     ██╔══╝  
╚██████╔╝██████╔╝███████║███████╗██║  ██║ ╚████╔╝ ██║  ██║██████╔╝███████╗███████╗
 ╚═════╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝
                                                                                  

```

# Complete Observability Stack


- Deploying a full Observability Solution and Event Streaming System

- `Loki_Tempo_Prom_Grafana/`
  - `Helm charts` and configs to Deploy `Tempo` , `Loki` , `Prometheus` , and `Grafana`
  - `Tempo` - Distributed Tracing Backend to receive Distributed Traces from Applications and Services
  - `Loki` - Distributed Logging Backend to receive Logs 
  - `Prometheus` - Distributed Metrics Backend to store Metrics for the Cluster and Services
  - `Grafana` - Visualization Software to visualize and query the Traces from Tempo , Logs from Loki , and Metrics from Prometheus


- `OTEL_Collector/`
  - `Helm charts` and configs to Deploy a Centralized Scalable `OpenTelemetry Collector`
  - Includes a DaemonSet that will be used to interact with deployments and services
  - A Deployment that handles `Cluster Metadata` and `Kubernetes Events`
  - Integrated with the Observability Stack


- <font color="#ffdeeb">`gokafka/`</font>

  - `Golang` application that utilizes all of Observability Stack and has deployments to deploy this to the Cluster
  - After deploying - we can see the `Traces`, `Metrics`, and `Logs` in `Grafana`

- <font color="#ffdeeb">`jvmkafka/`</font>
  - A more advanced system that uses `AutomaticOpenTelemetrySDK` Configuration with a `Kafka Broker` to perform functionality such as create events using Producers and uses `Virtual Threads` to consume events and process them.

