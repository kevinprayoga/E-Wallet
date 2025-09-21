terraform {
  cloud {
    organization = "abdinegara-org"
    workspaces {
      name = "ewallet-infra"
    }
  }
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {}
