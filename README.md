# Full-Stack Microservices Kubernetes Project

This repository contains a complete, production-ready Kubernetes orchestration for a microservices application, featuring a **Go Backend**, **React Frontend**, and **PostgreSQL Database**. The project demonstrates advanced DevOps practices including GitOps, automated CI/CD, and a full observability stack.

---

## Project Structure

* **`src/`**: Source code for the Frontend and Backend services.
* **`charts/`**: Helm charts for individual services (Backend, Frontend, Database, Loki, Jaeger).
* **`cd-root/`**: The "App-of-Apps" root chart for ArgoCD.
* **`tasks/task-x/`**: Dedicated folders for each assessment task (1–5) containing specific manifests and documentation.
* **`Makefile`**: The central command center for cluster management and deployment.

---

## Tech Stack

* **Orchestration**: Kubernetes (Kind)
* **CI/CD**: GitHub Actions
* **GitOps**: ArgoCD
* **Networking**: NGINX Ingress Controller & Istio Service Mesh
* **Observability**: Prometheus, Grafana, Loki, and Jaeger
* **Security**: Sealed Secrets, RBAC, and Image Scanning

---

## Quick Start

Ensure you have `docker`, `kind`, `kubectl`, and `helm` installed.

### 1. Complete Deployment

Run the following command to create the cluster, install the Ingress controller, setup ArgoCD, and deploy all applications:

```bash
make deploy-all

```

### 2. Accessing the Services

Once the deployment finishes, use these commands in separate terminals to access the stack:

| Service | Command | URL |
| --- | --- | --- |
| **App (Ingress)** | `make pf-ingress` | [http://localhost:8084](https://www.google.com/search?q=http://localhost:8084) |
| **Grafana** | `make pf-grafana` | [http://localhost:8081](https://www.google.com/search?q=http://localhost:8081) |
| **ArgoCD UI** | `kubectl port-forward -n argocd svc/argocd-server 8083:443` | [https://localhost:8083](https://www.google.com/search?q=https://localhost:8083) |

*To get your ArgoCD password, run:* `make argocd-password`

---

## Security & Reliability

* **Zero-Trust Networking**: Implemented via Kubernetes Network Policies and Istio mTLS.
* **Self-Healing**: Configured with Liveness and Readiness probes to ensure traffic only hits healthy pods.
* **HPA**: The backend scales automatically based on CPU load.
* **Scanning**: Every image is scanned by **Trivy** during the CI process to block vulnerabilities.

---

## Observability

* **Logging**: Centralized logs via **Loki**, viewable in Grafana.
* **Tracing**: Distributed tracing via **Jaeger** to debug request latency across services.
* **Metrics**: Custom application metrics scraped by **Prometheus**.

---

## Cleanup

To delete the cluster and all resources:

```bash
make down

```
