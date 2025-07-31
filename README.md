
## Folder Structure

google-checker/
├── main.go               # The Go app
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Dockerfile            # Docker build
├── k8s/
│   ├── deployment.yaml   # Kubernetes Deployment
│   ├── service.yaml      # Kubernetes Service
│   ├── prometheus.yaml   # Prometheus config (optional in-cluster)
│   └── grafana.yaml      # Grafana config (optional in-cluster)
└── README.md             # (Optional) setup notes

## 🔍 What Each File Does
main.go: Contains the checkGoogle() function and Prometheus /metrics endpoint.

go.mod / go.sum: Tracks dependencies (e.g. Prometheus client).

Dockerfile: Builds the Go app into a small, production-ready container.

k8s/deployment.yaml: Deploys the container as a pod with port 8080.

k8s/service.yaml: Exposes the app inside the cluster (ClusterIP or NodePort).

k8s/prometheus.yaml: Optional setup for deploying Prometheus that scrapes your service.

k8s/grafana.yaml: Optional setup for deploying Grafana and connecting it to Prometheus.

## ✅ Go Project Initialization Steps
Make sure you’re inside your project folder:

```bash
cd google-checker/
```

🥇 Step 1: Initialize Go module

```bash
go mod init github.com/yourname/google-checker
```

Replace yourname with your GitHub username or any unique prefix you like. It doesn't need to be real unless you plan to publish it.

🥈 Step 2: Add Dependencies

In your main.go, you'll use:

prometheus/client_golang

the standard net/http, log, etc.

After writing your basic main.go, run:

```bash
go mod tidy
```

This adds the necessary modules to go.mod and locks versions in go.sum.

✅ Result
You’ll now have:

```python
go.mod    # contains your module name + dependencies
go.sum    # dependency checksums
```

🔍 Want to test it?
After running go mod tidy, try:

```bash
go run main.go
```
Then open: http://localhost:8080/metrics

## Build and run locally to test:

```bash
docker build -t yourname/google-checker .
docker run -p 8080:8080 yourname/google-checker
```

Visit http://localhost:8080/metrics — do you see google_up 1?

## Build Docker image for minikube

minikube start

eval $(minikube docker-env)

Build Docker Image for Minikube
Run this command before building the image, so it builds inside Minikube’s Docker daemon:

```bash
eval $(minikube docker-env)
```Then build your Go app:```

```bash
docker build -t google-checker:local .
```

You can now use google-checker:local directly in your K8s manifests without pushing it anywhere.

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

```bash
minikube service google-checker --url
```

## Install Prometheus with Helm
✅ Step 1: Add the Prometheus Helm repo

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

✅ Step 2: Install Prometheus

```bash
helm install prometheus prometheus-community/prometheus \
  --namespace monitoring --create-namespace
```

This creates:

A new monitoring namespace

Prometheus Server

Node Exporter, Pushgateway, etc.

A ServiceMonitor CRD if needed (not used by default config, but handy with Prometheus Operator)

Step 3: Expose Prometheus (optional, for testing)
To view the Prometheus UI:

```bash
kubectl port-forward -n monitoring svc/prometheus-server 9090:80
```

Then open: http://localhost:9090

Step 4: Configure Prometheus to scrape your service
By default, Prometheus doesn’t auto-discover custom services in other namespaces unless you tell it how.

👇 Option A (simple): Annotate your google-checker pod
Add these annotations to your Deployment’s pod template:

```yaml
    metadata:
      labels:
        app: google-checker
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/path: /metrics
        prometheus.io/port: "8080"
```

Then re-apply your Deployment:

```bash
kubectl apply -f k8s/deployment.yaml
```

Set Up Grafana to Visualize It

Step 1: Install Grafana via Helm

```bash
helm install grafana grafana/grafana \
  --namespace monitoring \
  --set adminPassword=admin \
  --set service.type=NodePort \
  --create-namespace
```

This:

Installs Grafana into the monitoring namespace

Exposes it via NodePort

Sets admin password to admin

✅ Step 2: Get Grafana URL
Run:

```bash
minikube service grafana -n monitoring --url
```

You’ll get something like:
```bash
http://127.0.0.1:31370
```

Login with:

User: admin

Password: admin

✅ Step 3: Add Prometheus as a data source
In the Grafana UI:

⚙️ Gear icon (left sidebar) → Data Sources

Click Add data source

Select Prometheus

Set URL to:

```pgsql
http://prometheus-server.monitoring.svc.cluster.local
```

Click Save & Test

✅ Step 4: Create a simple dashboard
➕ Create → Dashboard

Add a panel

In the query box, type:

```bash
google_up
```

You should see values like 1 or 0

Customize visualization → Gauge, Time series, etc.