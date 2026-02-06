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
  description             = "Secret managed by arcctl: ${var.name}"
  recovery_window_in_days = 0

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id     = aws_secretsmanager_secret.this.id
  secret_string = var.data
}
