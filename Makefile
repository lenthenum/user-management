# --- Variables ---
CLUSTER_NAME=test
BACKEND_NAMESPACE=backend
FRONTEND_NAMESPACE=frontend
FRONTEND_IMG=lenthenum/um-frontend:0.0.1
BACKEND_IMG=lenthenum/um-backend:0.0.1
UI_PORT=3001
API_PORT=8002
DB_PORT=5432

# --- Cluster Management ---
cluster-up:
	kind create cluster --name $(CLUSTER_NAME)
	kind export kubeconfig --name $(CLUSTER_NAME)

cluster-down:
	kind delete cluster --name $(CLUSTER_NAME)

# --- Build & Load ---
build-frontend:
	docker build -t $(FRONTEND_IMG) ./src/frontend
	kind load docker-image $(FRONTEND_IMG) --name $(CLUSTER_NAME)

build-backend:
	docker build -t $(BACKEND_IMG) ./src/backend
	kind load docker-image $(BACKEND_IMG) --name $(CLUSTER_NAME)

# --- ArgoCD Operations ---
.PHONY: deploy-argocd
deploy-argocd:
	kubectl create namespace argocd
	kubectl apply -n argocd --server-side --force-conflicts -f https://raw.githubusercontent.com/argoproj/argo-cd/v3.3.2/manifests/install.yaml
	kubectl -n argocd wait --for=condition=available --timeout=300s --all deployments

.PHONY: argocd-password
argocd-password:
	@kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

.PHONY: login-argocd
login-argocd:
	kubectl port-forward service/argocd-server -n argocd 8083:443 &          
	argocd login localhost:8083 --grpc-web --insecure --username admin --password $$(kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

# --- Troubleshooting & Access ---
pf-ui:
	kubectl port-forward svc/um-frontend-service $(UI_PORT):3000 -n $(FRONTEND_NAMESPACE)

pf-api:
	kubectl port-forward svc/um-backend-service $(API_PORT):8000 -n $(BACKEND_NAMESPACE)

pf-db:
	kubectl port-forward svc/postgres-service $(DB_PORT):5432 -n $(BACKEND_NAMESPACE)

logs-db:
	kubectl logs -l app=postgres -n $(BACKEND_NAMESPACE) -f

# --- The "Total Reset" Button ---
reset-db:
	helm uninstall database -n $(BACKEND_NAMESPACE)
	kubectl delete pvc --all -n $(BACKEND_NAMESPACE)
	helm install database ./charts/database -n $(BACKEND_NAMESPACE)