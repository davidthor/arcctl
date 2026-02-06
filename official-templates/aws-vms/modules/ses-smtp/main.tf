terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

# SES SMTP configuration
# The actual SES domain identity and SMTP credentials are managed at the
# datacenter level via variables. This module serves as a placeholder
# to validate that SES is properly configured.

data "aws_ses_domain_identity" "this" {
  count  = var.identity_arn != "" ? 1 : 0
  domain = element(split("/", var.identity_arn), length(split("/", var.identity_arn)) - 1)
}
