# Task 2: CI/CD & GitOps Implementation

## 1. Overview

This task establishes a fully automated **Continuous Integration (CI)** and **Continuous Delivery (CD)** pipeline. It transitions the project from manual `kubectl` updates to a **GitOps** model using GitHub Actions and ArgoCD.

## 2. CI/CD Pipeline Design

### - Continuous Integration (GitHub Actions)

* **Linting & Analysis**: Automated linting ensures code quality and adherence to standards before testing begins.
* **Testing**: Unit and integration tests are executed on every Pull Request to prevent regressions.
* **Security Scanning**: A vulnerability scan is performed on the container image after the build but before the push. If "CRITICAL" vulnerabilities are detected, the pipeline fails and blocks the deployment.

### - Continuous Delivery & GitOps (ArgoCD)

* **Automated Versioning**: Every merge to `main` triggers a semantic version bump. The CI automatically updates the `Chart.yaml` in the repository.
* **ArgoCD Sync**: ArgoCD monitors the Helm charts. When it detects a version bump in git, it automatically pulls the new image and performs a rolling update in the cluster.
* **Separation of Concerns**: The CI system has no direct access to the Kubernetes cluster. It only has permission to write to the Git repository, significantly reducing the attack surface.

## 3. Review & Rollback Strategy

### - Branch Protection

* **Rule**: All changes to `main` must come through a Pull Request.
* **Requirement**: Status checks (CI tests, Security Scans) must pass, and at least one peer review is required.

### - Rollback

* **Git Revert**: Since Git is the "Source of Truth," reverting a commit to a previous version of `Chart.yaml` triggers ArgoCD to restore the previous stable state of the cluster.
* **ArgoCD Manual Rollback**: For emergency situations, the ArgoCD UI allows for an instant "Rollback" to any previous successful synchronization point.

### - Promotion

* Within ArgoCD application, you can override values to suit an environment (staging / production) with multiple values files.
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: um-backend-prod
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/lenthenum/um-stack.git
    targetRevision: HEAD
    path: charts/backend
    helm:
      # Override staging values for production
      valueFiles:
        - values.yaml
        - values-prod.yaml 
```

## 4. Further Improvements

1. Currently one github action workflow has single job which makes execution serial. For faster execution, split the action into multiple jobs.
2. Add an ArgoCD application to control ArgoCD itself.
3. Implement multiple values file structure to override values per environment type, hyperscalar and landscape
