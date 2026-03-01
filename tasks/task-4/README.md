### 5. README.md: Approach and Design Decisions

# Task 4: Security Hardening

## 1. Infrastructure Security

* **RBAC Implementation**: Defined three distinct personas (Admin, Operator, Developer). Use `RoleBindings` to restrict developers to specific namespaces (`frontend`/`backend`), preventing cross-tenant interference.
* **Pod Security Standards (PSS)**: Enforce `restricted` profile across all application namespaces for a **read-only root filesystem** and preventing privilege escalation.
Example -
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: backend
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted

```

* Configure security contexts at pod and container level for runAsNonRoot: true and readOnlyRootFilesystem: true wherever possible

## 2. Data & Secret Security

* **Sealed Secrets**: Implemented Bitnami Sealed Secrets to adhere to GitOps best practices. Secrets are encrypted at the client-side and can only be decrypted by the cluster-resident controller, ensuring sensitive data like `DATABASE_URL` is never exposed in plain text in the repository.
* **Network Segmentation**: Implement a "Default Deny" ingress policy, only allowing authorized paths between the Frontend, Backend, and Database.

## 3. Security in CI

* **Image Scanning**: Integrated **Trivy** into the GitHub Actions CI/CD pipeline. Images are scanned for CVEs at the build stage, and the pipeline is configured to fail if "CRITICAL" vulnerabilities are detected.


## 4. mTLS

Add **mTLS** (Mutual TLS) using a Service Mesh like **Istio** or **Linkerd**.