variable "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "Kubernetes context to target"
  type        = string
  default     = "k3d-devops-cluster"
}

variable "argocd_namespace" {
  description = "Namespace for ArgoCD installation"
  type        = string
  default     = "argocd"
}

variable "argo_rollouts_namespace" {
  description = "Namespace for Argo Rollouts installation"
  type        = string
  default     = "argo-rollouts"
}

variable "monitoring_namespace" {
  description = "Namespace for Prometheus & Grafana stack"
  type        = string
  default     = "monitoring"
}
