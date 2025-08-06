
## Info

I was refreshing my memory on some Kubernetes concepts including exporting Prometheus metrics from a container.  Noted that Chatgpt had a new Study and Learn Tool as of 2025-07-31 and decided to see where that took me.

## Folder Structure

```google-checker/
├── main.go               # The Go app
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Dockerfile            # Docker build
├── k8s/
│   ├── deployment.yaml   # Kubernetes Deployment
│   ├── service.yaml      # Kubernetes Service
│   ├── prometheus.yaml   # Prometheus config (optional in-cluster)
│   └── grafana.yaml      # Grafana config (optional in-cluster)
└── README.md             # (Optional) setup notes```

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

## Set Up Grafana to Visualize It

### Step 1: Install Grafana via Helm

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```


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

![alt text](image-1.png)



✅ Step 4: Create a simple dashboard
➕ Create → Dashboard

Add a panel

In the query box, type:

```bash
google_up
```

You should see values like 1 or 0

Customize visualization → Gauge, Time series, etc.

![alt text](image.png)

Absolutely — here are the **Slack alerting steps for Grafana** in clean Markdown format:

---

Here you go — clean and complete **Markdown version** of the steps to create an alert rule in Grafana:

---

````markdown
# 🔔 How to Create an Alert Rule in Grafana

## ✅ 1. Create or Edit a Panel

1. Open Grafana
2. Go to an existing dashboard or create a new one:
   - Click **+ → Dashboard → Add new panel**
3. Set the **Query** to:
   ```promql
   google_up
````

4. Choose a visualization type like **Stat** or **Gauge**
5. Click **Apply** to save the panel

---

## ✅ 2. Add an Alert Rule to the Panel

1. Hover over the panel → click the ⚙️ **Edit** button
2. Go to the **Alert** tab
3. Click **Create alert rule**

---

## ✅ 3. Define the Alert Rule

| Field            | Value                                  |
| ---------------- | -------------------------------------- |
| Rule name        | `Google Down Alert`                    |
| Folder           | `google-checker` or `Default`          |
| Evaluation group | `google_alerts` or `every_30s`         |
| Evaluate every   | `30s`                                  |
| Condition        | WHEN `google_up` IS BELOW `1` FOR `1m` |

---

## ✅ 4. Add Labels and Summary (Optional)

| Field   | Example                                      |
| ------- | -------------------------------------------- |
| Labels  | `severity=critical`, `source=google-checker` |
| Summary | `Google is unreachable`                      |

---

## ✅ 5. Save the Rule

1. Click **Save rule**
2. Alert rule will start evaluating immediately

---

## ✅ 6. Test the Alert

To simulate a failure:

```bash
kubectl scale deploy google-checker --replicas=0
```

Then go to **Alerting → Alert rules** in Grafana and look for:

* `Pending` → condition true, waiting to fire
* `Firing` → alert triggered
* `OK` → alert resolved

Scale it back up with:

```bash
kubectl scale deploy google-checker --replicas=1
```

---


# 📣 Setting Up Slack Alerts in Grafana

## ✅ 1. Create a Slack Webhook

1. Go to: [https://api.slack.com/apps](https://api.slack.com/apps)
2. Click **Create New App**
3. Choose **From scratch** and name it (e.g., `Grafana Alerts`)
4. Select your Slack workspace
5. In the left sidebar, click **Incoming Webhooks**
6. Toggle **Activate Incoming Webhooks** to **On**
7. Scroll down and click **Add New Webhook to Workspace**
8. Choose a channel (e.g., `#alerts`) and click **Allow**
9. Copy the generated **Webhook URL**, e.g.:

   ```
   https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
   ```

---

## ✅ 2. Add Slack as a Contact Point in Grafana

1. In Grafana UI, go to **Alerting → Contact points**
2. Click **New contact point**
3. Name: `Slack Alerts`
4. Type: **Slack**
5. Paste your Slack Webhook URL
6. (Optional) Customize the message or alert format
7. Click **Save contact point**

---

## ✅ 3. Create a Notification Policy

1. Go to **Alerting → Notification policies**
2. Click **New policy** or edit the existing one
3. Leave **Matching condition** blank to match all alerts
4. Set **Send to** → `Slack Alerts`
5. Click **Save policy**

---

## ✅ 4. Test the Alert

You can test alerts two ways:

### Option A: Test manually from Grafana

* Go to **Alerting → Contact points**
* Click **Test** on your Slack contact point

### Option B: Trigger a real alert

* Scale down your app to simulate failure:

  ```bash
  kubectl scale deploy google-checker --replicas=0
  ```

* Wait 1–2 minutes for `google_up == 0`

* Grafana should fire an alert and send a Slack message

---

## ✅ 5. Scale the app back up

```bash
kubectl scale deploy google-checker --replicas=1
```

---

Let me know if you want a **prebuilt alert rule template** or want to tweak the Slack message formatting.


## Persistence

## Tracing

Add Jaeger 

```bash
kubectl apply -f k8s/jaeger.yaml
```

Expose Jaeger Web UI via Minikube:

```bash
minikube service jaeger -n observability
```

Open Jaeger Web UI, for example http://127.0.0.1:33117

## OpenTelemetry

Add modules

```bash
go get go.opentelemetry.io/otel@latest
go get go.opentelemetry.io/otel/sdk@latest
go get go.opentelemetry.io/otel/exporters/jaeger@latest
```

