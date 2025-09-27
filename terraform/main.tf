# container registry
resource "digitalocean_container_registry" "ewallet-registry" {
  name                    = var.registry_name
  subscription_tier_slug  = var.registry_tier
  region                  = var.region
}

# kubernetes cluster
resource "digitalocean_kubernetes_cluster" "ewallet-cluster" {
  name    = var.cluster_name
  region  = var.region
  version = var.k8s_version

  registry_integration = true

  node_pool {
    name       = "worker-node-pool"
    size       = var.node_size
    node_count = var.node_count
    auto_scale = false
  }
}

# ambil docker credentials untuk login registry
resource "digitalocean_container_registry_docker_credentials" "creds" {
  registry_name  = digitalocean_container_registry.ewallet-registry.name
  expiry_seconds = 86400
  write          = true
}
