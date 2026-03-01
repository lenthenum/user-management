# --- Variables ---
CLUSTER_NAME       ?= um-cluster
BACKEND_NAMESPACE  ?= backend
FRONTEND_NAMESPACE ?= frontend
INGRESS_NS         ?= ingress-nginx
GRAFANA_NS         ?= grafana

# --- The "Single Command" Deployer ---
.PHONY: deploy-all
deploy-all: up deploy-ingress deploy-argocd deploy-app
	@echo "-------------------------------------------------------"
	@echo "✅ All systems deployed!"
	@echo "🔐 ArgoCD Password: $$(make -s argocd-password)"
	@echo "🌐 Ingress Entry: http://localhost:8084"
	@echo "📊 Grafana Dashboards: http://localhost:8081"
	@echo "-------------------------------------------------------"

# --- Existing Management Targets ---
.PHONY: up
up:
	kind create cluster --name $(CLUSTER_NAME) --config kind-config.yaml || true
	kind export kubeconfig --name $(CLUSTER_NAME)

.PHONY: deploy-ingress
deploy-ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	@echo "Waiting for Ingress controller..."
	kubectl wait --namespace $(INGRESS_NS) --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s

.PHONY: deploy-argocd
deploy-argocd:
	kubectl create namespace argocd || true
	kubectl apply -n argocd --server-side --force-conflicts -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
	kubectl -n argocd wait --for=condition=available --timeout=300s deployment/argocd-server

.PHONY: deploy-app
deploy-app:
	helm upgrade --install root ./cd-root -n argocd --create-namespace

# --- Granular Port-Forwarding Targets ---

.PHONY: pf-ingress
pf-ingress:
	@echo "Opening Ingress Gateway on http://localhost:8084"
	kubectl port-forward svc/ingress-nginx-controller -n $(INGRESS_NS) 8084:80

.PHONY: pf-grafana
pf-grafana:
	@echo "Opening Grafana on http://localhost:8081"
	kubectl port-forward svc/loki-stack-app-grafana -n $(GRAFANA_NS) 8081:80

.PHONY: pf-apps
pf-apps:
	@echo "Opening Apps (UI:3001, API:8002, DB:5432)"
	kubectl port-forward svc/um-frontend-service 3001:80 -n $(FRONTEND_NAMESPACE) &
	kubectl port-forward svc/um-backend-service 8002:8000 -n $(BACKEND_NAMESPACE) &
	kubectl port-forward svc/postgres-service 5432:5432 -n $(BACKEND_NAMESPACE) &

.PHONY: pf-all
pf-all:
	@make -j 3 pf-ingress pf-grafana pf-apps

# --- Helper Targets ---
.PHONY: argocd-password
argocd-password:
	@kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo

.PHONY: down
down:
	kind delete cluster --name $(CLUSTER_NAME)