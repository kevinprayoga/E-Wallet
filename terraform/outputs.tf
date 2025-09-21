output "registry_endpoint" {
  description = "Endpoint Container Registry"
  value       = digitalocean_container_registry.ewallet-registry.endpoint
}

output "kubeconfig" {
  description = "Kubeconfig untuk mengakses cluster"
  value       = digitalocean_kubernetes_cluster.ewallet-cluster.kube_config[0].raw_config
  sensitive   = true
}

output "docker_credentials" {
  description = "Base64 Docker config.json untuk login registry"
  value       = digitalocean_container_registry_docker_credentials.creds.docker_credentials
  sensitive   = true
}
