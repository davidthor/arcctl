terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
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

# --- Asymmetric keys (RSA / ECDSA) ---

resource "tls_private_key" "asymmetric" {
  count     = var.key_type != "symmetric" ? 1 : 0
  algorithm = upper(var.key_type)
  rsa_bits  = var.key_type == "rsa" ? var.key_size : null
}

# --- Symmetric keys ---

resource "random_bytes" "symmetric" {
  count  = var.key_type == "symmetric" ? 1 : 0
  length = var.key_size / 8
}
