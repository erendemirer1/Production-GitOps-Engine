# ==============================================================================
# 1. ArgoCD GitOps Engine
# ==============================================================================
resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = "6.7.18"
  namespace        = var.argocd_namespace
  create_namespace = true

  set {
    name  = "server.service.type"
    value = "ClusterIP"
  }
}

# ==============================================================================
# 2. Argo Rollouts (Progressive Canary Delivery)
# ==============================================================================
resource "helm_release" "argo_rollouts" {
  name             = "argo-rollouts"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-rollouts"
  version          = "2.36.0"
  namespace        = var.argo_rollouts_namespace
  create_namespace = true

  set {
    name  = "dashboard.enabled"
    value = "true"
  }
}

# ==============================================================================
# 3. Kube-Prometheus-Stack (Prometheus, Alertmanager, Grafana)
# ==============================================================================
resource "helm_release" "prometheus_stack" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = "58.2.2"
  namespace        = var.monitoring_namespace
  create_namespace = true

  set {
    name  = "grafana.adminPassword"
    value = "admin"
  }

  set {
    name  = "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues"
    value = "false"
  }
}
