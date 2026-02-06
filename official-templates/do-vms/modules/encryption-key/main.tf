terraform {
  required_providers {
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

locals {
  is_asymmetric = var.key_type == "rsa" || var.key_type == "ecdsa"
  is_symmetric  = var.key_type == "symmetric"
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
