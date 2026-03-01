### Task 5: Istio Service Mesh Analysis

**1. Role of Istio and the Sidecar Model**

* Istio acts as a "service mesh" layer that manages how microservices talk to each other, handling security and traffic without us needing to change our application code.
* It automatically adds an envoy proxy container (sidecar) to every pod which intercepts all network traffic entering or leaving the application.
* It moves networking logic—like retries, timeouts, and encryption—from the app code to the infrastructure, ensuring these features work the same way for every service.

**2. PeerAuthentication vs. AuthorizationPolicy**

* **PeerAuthentication:** This controls the "transport" security which is if the connection between two services is encrypted using mTLS.
* **AuthorizationPolicy:** This handles "access" security. It checks if a specific service is actually allowed to talk to another based on its identity or the path it is trying to access.
* **Strict mTLS:** To enforce this cluster-wide or by namespace, you apply a `PeerAuthentication` policy with the mode set to `STRICT`, which forces all pods to reject any unencrypted traffic.

**3. Traffic Management and Canary Deployments**

* Istio sends routing rules to all proxies, allowing us to control traffic flow dynamically.
* **Canary Walkthrough:**
* **Step 1:** Create a DestinationRule to define different versions (subsets) of the service based on pod labels (e.g., `v1` and `v2`).
* **Step 2:** Configure a VirtualService to split traffic by weight, such as sending 90% of requests to the stable version v1 and 10% to the new canary version v2.

**4. Istio Ingress Gateway vs. Standard Ingress**

* **Istio Ingress Gateway:** This is a dedicated Envoy proxy sitting at the edge of the cluster; it integrates fully with the rest of the mesh, allowing for consistent security and routing policies from the very first entry point.
* **Standard Ingress:** Usually a basic controller (like NGINX) that handles external routing but lacks the deep "mesh" features like automatic mTLS or advanced telemetry that the Istio gateway provides.

**5. Observability Integration**

* As all traffic passes through the sidecar proxies, Istio can collect data on every request without requiring any extra code in the apps.
* Istio sends data like request rates and error counts to **Prometheus**, which is then visualized in **Grafana** dashboards.
* Istio passes along "trace headers" between services, allowing **Jaeger** to show the full path of a request through the system.
