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
  is_encryption_key = var.key_type != null
  is_asymmetric     = var.key_type == "rsa" || var.key_type == "ecdsa"
  is_symmetric      = var.key_type == "symmetric"
  is_data_secret    = var.data != null
}

# RSA key generation
resource "tls_private_key" "rsa" {
  count     = var.key_type == "rsa" ? 1 : 0
  algorithm = "RSA"
  rsa_bits  = coalesce(var.key_size, 2048)
}

# ECDSA key generation
resource "tls_private_key" "ecdsa" {
  count       = var.key_type == "ecdsa" ? 1 : 0
  algorithm   = "ECDSA"
  ecdsa_curve = "P${coalesce(var.key_size, 256)}"
}

# Symmetric key generation
resource "random_bytes" "symmetric" {
  count  = local.is_symmetric ? 1 : 0
  length = coalesce(var.key_size, 256) / 8
}

# For plain data secrets, generate a unique ID
resource "random_id" "secret_id" {
  count       = local.is_data_secret && !local.is_encryption_key ? 1 : 0
  byte_length = 8
  prefix      = "${var.name}-"
}
