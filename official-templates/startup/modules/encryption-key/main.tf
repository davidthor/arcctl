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

# RSA key pair (when key_type == "rsa")
resource "tls_private_key" "rsa" {
  count = var.key_type == "rsa" ? 1 : 0

  algorithm = "RSA"
  rsa_bits  = var.key_size != null ? var.key_size : 4096
}

# ECDSA key pair (when key_type == "ecdsa")
resource "tls_private_key" "ecdsa" {
  count = var.key_type == "ecdsa" ? 1 : 0

  algorithm   = "ECDSA"
  ecdsa_curve = var.key_size != null ? "P${var.key_size}" : "P256"
}

# Symmetric key (when key_type == "symmetric")
resource "random_bytes" "symmetric" {
  count = var.key_type == "symmetric" ? 1 : 0

  length = var.key_size != null ? var.key_size / 8 : 32
}

locals {
  # Select the appropriate key based on type
  is_rsa   = var.key_type == "rsa"
  is_ecdsa = var.key_type == "ecdsa"
  is_sym   = var.key_type == "symmetric"

  private_key_pem = local.is_rsa ? try(tls_private_key.rsa[0].private_key_pem, "") : (
    local.is_ecdsa ? try(tls_private_key.ecdsa[0].private_key_pem, "") : ""
  )

  public_key_pem = local.is_rsa ? try(tls_private_key.rsa[0].public_key_pem, "") : (
    local.is_ecdsa ? try(tls_private_key.ecdsa[0].public_key_pem, "") : ""
  )

  symmetric_key = local.is_sym ? try(random_bytes.symmetric[0].hex, "") : ""
}
