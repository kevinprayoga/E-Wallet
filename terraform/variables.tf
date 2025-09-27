variable "DIGITALOCEAN_TOKEN" {
  description = "Token API DigitalOcean"
  type        = string
  sensitive   = true
}

variable "registry_name" {
  description = "Nama DigitalOcean Container Registry"
  type = string
}

variable "registry_tier" {
  type = string
  default = "starter"
}

variable "region" {
  description = "Region DigitalOcean"
  type    = string
}

variable "cluster_name" {
  description = "Nama Kubernetes Cluster"
  type    = string
}

variable "k8s_version" {
  type    = string
  default = "1.33.1-do.4"
}

variable "node_size" {
  description = "Ukuran node worker"
  type        = string
}

variable "node_count" {
  description = "Jumlah node worker"
  type        = number
}
