terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "random_id" "suffix" {
  byte_length = 4
}

locals {
  repo_name = "arcctl/${replace(basename(var.context), ".", "app")}-${random_id.suffix.hex}"
  image_tag = "latest"
  image_uri = "${aws_ecr_repository.this.repository_url}:${local.image_tag}"

  dockerfile_arg = var.dockerfile != null ? "-f ${var.dockerfile}" : ""
  target_arg     = var.target != null ? "--target ${var.target}" : ""
  build_args     = var.args != null ? join(" ", [for k, v in var.args : "--build-arg ${k}=${v}"]) : ""
}

resource "aws_ecr_repository" "this" {
  name                 = local.repo_name
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name      = local.repo_name
    ManagedBy = "arcctl"
  }
}

resource "null_resource" "docker_build" {
  triggers = {
    context    = var.context
    dockerfile = var.dockerfile != null ? var.dockerfile : ""
    target     = var.target != null ? var.target : ""
    args       = var.args != null ? jsonencode(var.args) : ""
    always     = timestamp()
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws ecr get-login-password --region ${data.aws_region.current.name} | \
        docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com

      docker build ${local.dockerfile_arg} ${local.target_arg} ${local.build_args} \
        -t ${local.image_uri} ${var.context}

      docker push ${local.image_uri}
    EOT
  }

  depends_on = [aws_ecr_repository.this]
}
