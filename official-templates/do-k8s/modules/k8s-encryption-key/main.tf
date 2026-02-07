terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "kubernetes" {
  config_path = local.kubeconfig_path
}

locals {
  kubeconfig_path = "${path.module}/.kubeconfig"
  is_asymmetric   = var.key_type == "rsa" || var.key_type == "ecdsa"
  is_symmetric    = var.key_type == "symmetric"
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

# RSA key generation
resource "tls_private_key" "rsa" {
  count     = var.key_type == "rsa" ? 1 : 0
  algorithm = "RSA"
  rsa_bits  = var.key_size
}

# ECDSA key generation
resource "tls_private_key" "ecdsa" {
  count       = var.key_type == "ecdsa" ? 1 : 0
  algorithm   = "ECDSA"
  ecdsa_curve = "P${var.key_size}"
}

# Symmetric key generation
resource "random_bytes" "symmetric" {
  count  = local.is_symmetric ? 1 : 0
  length = var.key_size / 8
}

# Store key material in Kubernetes secret
resource "kubernetes_secret_v1" "key_secret" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/managed-by" = "cldctl"
      "cldctl/resource-type"         = "encryption-key"
      "cldctl/key-type"              = var.key_type
    }
  }

  data = local.is_asymmetric ? {
    "private_key" = var.key_type == "rsa" ? tls_private_key.rsa[0].private_key_pem : tls_private_key.ecdsa[0].private_key_pem
    "public_key"  = var.key_type == "rsa" ? tls_private_key.rsa[0].public_key_pem : tls_private_key.ecdsa[0].public_key_pem
  } : {
    "key" = random_bytes.symmetric[0].hex
  }

  type = "Opaque"

  depends_on = [local_file.kubeconfig]
}
