# Capp User Guide

## Introduction

**Capp** (Container Application) is a Kubernetes Custom Resource that provides a simplified abstraction for deploying containerized serverless workloads. It allows users to deploy applications without requiring deep knowledge of Kubernetes concepts, while automatically managing autoscaling, routing, logging, and storage. For architecture details and operator installation instructions, please refer to the [main README](../README.md).

## Capp CR Specification

A Capp Custom Resource consists of several key fields that control different aspects of your application deployment. Below is a detailed explanation of each field and its functional requirements.

### `scaleMetric`

**Purpose**: Defines which metric the autoscaler uses to scale your application up or down.

**Values**:
- `concurrency` (default): Scales based on the number of concurrent requests being processed. Best for HTTP services with varied request durations.
- `rps`: Scales based on requests per second. Ideal for high-throughput APIs with predictable request patterns.
- `cpu`: Scales based on CPU utilization percentage. Suitable for CPU-intensive workloads.
- `memory`: Scales based on memory utilization percentage. Best for memory-intensive applications.

**Functional Requirement**: Must be one of the four supported values. The operator creates the appropriate autoscaler (HPA or KPA) based on this selection.

### `state`

**Purpose**: Controls whether the application workload is running or suspended.

**Values**:
- `enabled` (default): The application is running and serving traffic.
- `disabled`: The application is suspended (scaled to zero) but configuration is preserved.

**Functional Requirement**: Use `disabled` to temporarily stop an application without deleting it. This is useful for cost savings during non-peak hours or maintenance windows.

### `configurationSpec`

**Purpose**: Defines the container specification for your application, including image, environment variables, and resource requirements.

**Structure**: Based on Knative's ConfigurationSpec, it contains a `template` with a `spec` that includes:
- `containers`: Array of container definitions
  - `name`: Container name
  - `image`: Container image (e.g., `ghcr.io/myorg/myapp:v1.0.0`)
  - `env`: Environment variables as key-value pairs
  - `resources`: CPU and memory requests/limits
  - `volumeMounts`: Mount points for volumes

**Functional Requirement**: At least one container must be specified with a valid image. The container configuration follows standard Kubernetes pod specifications.

### `routeSpec`

**Purpose**: Configures custom DNS routing and TLS for accessing your application.

**Fields**:
- `hostname`: Custom DNS name for your application (e.g., `myapp.example.com`)
- `tlsEnabled`: Boolean to enable/disable HTTPS with automatic certificate management
- `trafficTarget`: Advanced traffic routing for canary deployments or A/B testing
- `routeTimeoutSeconds`: Maximum duration for a request before timeout

**Functional Requirements**:
- If `hostname` is specified, the operator creates a DomainMapping, DNS Record, and optionally a Certificate.
- TLS certificates are automatically provisioned when `tlsEnabled: true`.

### `logSpec`

**Purpose**: Configures automatic log shipping to Elasticsearch for centralized logging.

**Fields**:
- `type`: Log destination type (currently only `elastic` is supported)
- `host`: Elasticsearch host address (IP or hostname)
- `index`: Elasticsearch index name where logs will be stored
- `user`: Username for Elasticsearch authentication
- `passwordSecret`: Name of the Kubernetes secret containing the Elasticsearch password

**Functional Requirement**: When configured, the operator creates Flow and Output resources to automatically collect logs from your application's stdout and ship them to Elasticsearch.

### `volumesSpec`

**Purpose**: Defines persistent storage volumes to be mounted in your application containers.

**Structure**:
- `nfsVolumes`: Array of NFS volume definitions
  - `name`: Volume name (must match the `volumeMounts` name in container spec)
  - `server`: NFS server hostname or IP address
  - `path`: Export path on the NFS server
  - `capacity`: Storage size (e.g., `200Gi`)

**Functional Requirements**:
- The nfspvc-operator must be installed in the cluster.
- Volume names must match those referenced in `volumeMounts` within the container spec.
- Knative must be configured to support persistent volumes.

### `sources`

**Purpose**: Configures event sources that trigger your application, enabling event-driven architectures.

**Structure**:
- `name`: Source name
- `type`: Source type (currently `Kafka` is supported)
- `bootstrapServers`: Array of Kafka broker addresses
- `topic`: Array of Kafka topics to consume from
- `kafkaAuth`: Authentication configuration
  - `username`: Kafka username
  - `passwordKey`: Reference to secret containing password

**Functional Requirement**: When configured, your application receives events from the specified Kafka topics, enabling serverless event processing.

## How to Use Capp

This section provides step-by-step instructions for common Capp usage scenarios. These instructions assume the container-app-operator is already installed in your cluster.

### Step 1: Create a Basic Capp

Create a simple Capp with just a container and default settings:

```yaml
apiVersion: rcs.dana.io/v1alpha1
kind: Capp
metadata:
  name: my-app
  namespace: my-namespace
spec:
  configurationSpec:
    template:
      spec:
        containers:
          - name: my-app
            image: ghcr.io/myorg/my-app:v1.0.0
  state: enabled
```

Apply it with: `kubectl apply -f my-app.yaml`

### Step 2: Configure Autoscaling

Choose the appropriate scaling metric for your workload:

**For high-traffic APIs**:
```yaml
spec:
  scaleMetric: rps
```

**For CPU-intensive workloads**:
```yaml
spec:
  scaleMetric: cpu
```

**For memory-intensive applications**:
```yaml
spec:
  scaleMetric: memory
```

**For concurrent request handling** (default):
```yaml
spec:
  scaleMetric: concurrency
```

### Step 3: Add a Custom Domain with TLS

To expose your application with a custom domain and HTTPS:

```yaml
spec:
  routeSpec:
    hostname: myapp.example.com
    tlsEnabled: true
```

**Note**: Ensure the CappConfig in the operator namespace has the correct DNS configuration.

### Step 4: Enable Elasticsearch Logging

To automatically ship logs to Elasticsearch:

```yaml
spec:
  logSpec:
    type: elastic
    host: elasticsearch.example.com
    index: my-app-logs
    user: elastic
    passwordSecret: es-password-secret
```

First, create the secret:
```bash
kubectl create secret generic es-password-secret \
  --from-literal=password='your-es-password' \
  -n my-namespace
```

### Step 5: Mount NFS Volumes

To add persistent storage to your application:

```yaml
spec:
  configurationSpec:
    template:
      spec:
        containers:
          - name: my-app
            image: ghcr.io/myorg/my-app:v1.0.0
            volumeMounts:
              - name: data-volume
                mountPath: /data
  volumesSpec:
    nfsVolumes:
      - name: data-volume
        server: nfs.example.com
        path: /exports/my-app-data
        capacity:
          storage: 100Gi
```

### Step 6: Connect Kafka Event Sources

To make your application event-driven:

```yaml
spec:
  sources:
    - name: kafka-events
      type: Kafka
      bootstrapServers:
        - kafka-broker-1:9092
        - kafka-broker-2:9092
      topic:
        - user-events
        - order-events
      kafkaAuth:
        username: kafka-user
        passwordKey:
          name: kafka-secret
          key: password
```

Create the Kafka secret:
```bash
kubectl create secret generic kafka-secret \
  --from-literal=password='your-kafka-password' \
  -n my-namespace
```

### Step 7: Manage Capp State

**To disable (suspend) an application**:
```bash
kubectl patch capp my-app -n my-namespace --type=merge -p '{"spec":{"state":"disabled"}}'
```

**To re-enable an application**:
```bash
kubectl patch capp my-app -n my-namespace --type=merge -p '{"spec":{"state":"enabled"}}'
```

### Step 8: Check Capp Status

**View basic status**:
```bash
kubectl get capp my-app -n my-namespace
```

**View detailed status**:
```bash
kubectl describe capp my-app -n my-namespace
```

**Check the status section** for:
- `knativeObjectStatus`: Underlying Knative service status
- `routeStatus`: Domain mapping and DNS record status
- `loggingStatus`: Logging configuration status
- `volumesStatus`: NFS volume status
- `conditions`: Overall health conditions

## Practical Examples

### Example 1: Simple Web Application with Custom Domain

This example deploys a web application with a custom domain, TLS, and RPS-based autoscaling:

```yaml
apiVersion: rcs.dana.io/v1alpha1
kind: Capp
metadata:
  name: web-app
  namespace: production
spec:
  scaleMetric: rps
  state: enabled
  configurationSpec:
    template:
      spec:
        containers:
          - name: web-app
            image: ghcr.io/mycompany/web-app:v2.1.0
            env:
              - name: ENVIRONMENT
                value: production
              - name: PORT
                value: "8080"
            resources:
              requests:
                memory: "256Mi"
                cpu: "200m"
              limits:
                memory: "512Mi"
                cpu: "500m"
  routeSpec:
    hostname: web.mycompany.com
    tlsEnabled: true
    routeTimeoutSeconds: 60
```

**What this does**:
- Deploys a containerized web application
- Automatically scales based on requests per second
- Accessible via `https://web.mycompany.com`
- Automatic TLS certificate management
- 60-second request timeout
- Resource limits prevent overuse

### Example 2: Advanced Event-Driven Application

This example shows a comprehensive setup with logging, persistent storage, and Kafka event processing:

```yaml
apiVersion: rcs.dana.io/v1alpha1
kind: Capp
metadata:
  name: event-processor
  namespace: analytics
spec:
  scaleMetric: cpu
  state: enabled
  configurationSpec:
    template:
      spec:
        containers:
          - name: processor
            image: ghcr.io/mycompany/event-processor:v1.5.0
            env:
              - name: DATA_DIR
                value: /data
              - name: BATCH_SIZE
                value: "1000"
            resources:
              requests:
                memory: "1Gi"
                cpu: "500m"
              limits:
                memory: "2Gi"
                cpu: "1000m"
            volumeMounts:
              - name: processed-data
                mountPath: /data
  routeSpec:
    hostname: processor.analytics.mycompany.com
    tlsEnabled: true
  volumesSpec:
    nfsVolumes:
      - name: processed-data
        server: nfs-storage.internal
        path: /exports/analytics/processed
        capacity:
          storage: 500Gi
  logSpec:
    type: elastic
    host: elasticsearch.monitoring.svc.cluster.local
    index: event-processor-logs
    user: analytics-user
    passwordSecret: es-analytics-secret
  sources:
    - name: user-events
      type: Kafka
      bootstrapServers:
        - kafka-1.internal:9092
        - kafka-2.internal:9092
        - kafka-3.internal:9092
      topic:
        - user-activity
        - user-transactions
      kafkaAuth:
        username: analytics-consumer
        passwordKey:
          name: kafka-analytics-secret
          key: password
```

**What this does**:
- Processes events from Kafka topics in real-time
- Scales based on CPU utilization
- Stores processed data on NFS persistent storage (500Gi)
- Ships all logs to Elasticsearch for monitoring
- Accessible via HTTPS for status/health checks
- Resource limits ensure stable performance

**Before applying, create the required secrets**:

```bash
# Elasticsearch secret
kubectl create secret generic es-analytics-secret \
  --from-literal=password='es-password' \
  -n analytics

# Kafka secret
kubectl create secret generic kafka-analytics-secret \
  --from-literal=password='kafka-password' \
  -n analytics
```

## Troubleshooting Tips

**Application not scaling**:
- Check if the CappConfig has the correct autoscale target values
- Verify the scaleMetric is appropriate for your workload
- Check the Knative autoscaler metrics

**Custom domain not working**:
- Verify CappConfig has correct DNS configuration
- Check the DomainMapping status: `kubectl get domainmapping -n <namespace>`
- Verify DNS records were created: `kubectl get cnamerecord -n <namespace>`

**Logs not appearing in Elasticsearch**:
- Verify Elasticsearch credentials in the secret
- Check the Flow and Output resources: `kubectl get syslogngflow,syslogngoutput -n <namespace>`
- Ensure the logging-operator is running

**Volume mount issues**:
- Verify Knative has persistent volume support enabled
- Check NFS server accessibility from the cluster
- Verify the NfsPvc resource status: `kubectl get nfspvc -n <namespace>`

For more detailed troubleshooting, installation instructions, and architecture information, refer to the [main README](../README.md).

