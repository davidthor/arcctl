terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

locals {
  is_rsa   = var.key_type == "rsa"
  is_ecdsa = var.key_type == "ecdsa"
  rsa_bits = var.key_size != null ? var.key_size : 4096

  ecdsa_curve = var.key_size != null ? (
    var.key_size <= 256 ? "P256" :
    var.key_size <= 384 ? "P384" : "P521"
  ) : "P256"
}

resource "tls_private_key" "this" {
  algorithm   = upper(var.key_type)
  rsa_bits    = local.is_rsa ? local.rsa_bits : null
  ecdsa_curve = local.is_ecdsa ? local.ecdsa_curve : null
}

resource "aws_secretsmanager_secret" "private_key" {
  name                    = "${var.name}/private-key"
  description             = "Private key for ${var.name}"
  recovery_window_in_days = 0

  tags = {
    Name      = "${var.name}-private-key"
    ManagedBy = "arcctl"
  }
}

resource "aws_secretsmanager_secret_version" "private_key" {
  secret_id     = aws_secretsmanager_secret.private_key.id
  secret_string = tls_private_key.this.private_key_pem
}

resource "aws_secretsmanager_secret" "public_key" {
  name                    = "${var.name}/public-key"
  description             = "Public key for ${var.name}"
  recovery_window_in_days = 0

  tags = {
    Name      = "${var.name}-public-key"
    ManagedBy = "arcctl"
  }
}

resource "aws_secretsmanager_secret_version" "public_key" {
  secret_id     = aws_secretsmanager_secret.public_key.id
  secret_string = tls_private_key.this.public_key_pem
}
