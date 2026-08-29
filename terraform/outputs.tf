output "argocd_server_service" {
  description = "ArgoCD Server Service DNS"
  value       = "argocd-server.${var.argocd_namespace}.svc.cluster.local"
}

output "grafana_service" {
  description = "Grafana Service DNS"
  value       = "prometheus-grafana.${var.monitoring_namespace}.svc.cluster.local"
}

output "platform_status" {
  description = "Status of the deployed Cloud-Native Infrastructure"
  value       = "GitOps Engine (ArgoCD), Argo Rollouts, and Prometheus Stack bootstrapped successfully via Terraform."
}
