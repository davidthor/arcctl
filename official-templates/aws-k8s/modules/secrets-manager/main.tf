terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_secretsmanager_secret" "this" {
  name                    = var.name
  description             = "Secret managed by cldctl: ${var.name}"
  recovery_window_in_days = 0

  tags = {
    Name      = var.name
    ManagedBy = "cldctl"
  }
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id     = aws_secretsmanager_secret.this.id
  secret_string = var.data
}
