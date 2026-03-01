# Task 1: Kubernetes Orchestration Fundamentals

## Project Overview

This task aims to cover the deployment and orchestration of a multi-tier microservices application that consists of a **Next.js Frontend**, a **Go Backend API**, and a **PostgreSQL Database**, all running within a local Kubernetes cluster (Kind).

### 1. Strategy

* **Namespace Isolation**: Services are split into `frontend` and `backend` namespaces as per respective microservices for better isolation and control over security aspects and deployments.
* **Unified Ingress Strategy**: The Ingress controller is placed in the `backend` namespace to allow direct resolution of the API endpoints, reducing 503 errors and DNS latency.

### 2. Scalability & Availability

* **Horizontal Pod Autoscaler (HPA)**: Implemented for the Go Backend to scale replicas based on CPU utilization (target 50%).
* **Pod Disruption Budgets (PDB)**: Configured a `minAvailable: 1` policy for the backend to ensure the API remains reachable during node drains or cluster upgrades.

### 3. Reliability & Self-Healing

* **Health Checks**:
* **Liveness Probes**: Automatically restarts the backend container if the main process hangs.
* **Readiness Probes**: Ensures the Backend is not added to the Ingress pool until the DB migrations are successfully completed.

### 4. Security

* **Network Policies**: The Database only accepts traffic from the Backend pods, and the Backend only accepts traffic from the Ingress Controller.
* **Secrets Management**: Database credentials are retrieved from Kubernetes Secrets (created from sealed secrets) rather than plain-text environment variables.

### 5. Persistence

* **StatefulSets**: The PostgreSQL database uses a StatefulSet with a `volumeClaimTemplate` to ensure data persists across pod restarts and maintains a stable network identity.

## 4. Deployment Instructions

### Prerequisites

* A running Kubernetes cluster (Kind, Minikube, or EKS/GKE).
* `kubectl` CLI installed.
* NGINX Ingress Controller installed in the cluster.

### Step-by-Step Execution
[As manifests are copied from helm chart at the root, please set the values in the manifest or install the helm chart from the root]

If values are set in manifests of task folder -
1. **Apply Namespaces**: `kubectl apply -f manifests/01-namespaces.yaml`
2. **Deploy Database**: `kubectl apply -f manifests/02-database/`
3. **Deploy Backend**: `kubectl apply -f manifests/03-backend/`
4. **Deploy Frontend**: `kubectl apply -f manifests/04-frontend/`
5. **Apply Networking**: `kubectl apply -f manifests/05-networking/`

From the repository root -
1. helm install db ./charts/database -n backend --create-namespace
2. helm install api ./charts/backend -n backend
3. helm install ui ./charts/frontend -n frontend --create-namespace

3. **Access the Application**:
* Port-forward the Ingress: `kubectl port-forward deployment/ingress-nginx-controller -n ingress-nginx 8084:80`
* Open your browser to: `http://localhost:8084`
