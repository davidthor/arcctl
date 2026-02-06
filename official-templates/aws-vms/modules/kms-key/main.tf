terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

data "aws_caller_identity" "current" {}

# Generate a random 256-bit symmetric key
resource "random_bytes" "key_material" {
  length = 32
}

resource "aws_kms_key" "this" {
  description             = "Symmetric encryption key for ${var.name}"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "EnableRootAccountAccess"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action   = "kms:*"
        Resource = "*"
      },
    ]
  })

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}

resource "aws_kms_alias" "this" {
  name          = "alias/arcctl-${replace(var.name, "/", "-")}"
  target_key_id = aws_kms_key.this.key_id
}
