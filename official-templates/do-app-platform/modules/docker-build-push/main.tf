terraform {
  required_providers {
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

locals {
  dockerfile = coalesce(var.dockerfile, "Dockerfile")
  tag        = "${var.registry}/${random_string.tag.result}:latest"

  build_args = var.args != null ? join(" ", [
    for key, value in var.args : "--build-arg ${key}=${value}"
  ]) : ""

  target_arg = var.target != null ? "--target ${var.target}" : ""
}

resource "random_string" "tag" {
  length  = 16
  special = false
  upper   = false
}

resource "null_resource" "docker_build" {
  triggers = {
    context    = var.context
    dockerfile = local.dockerfile
    target     = var.target
    args       = jsonencode(var.args)
  }

  provisioner "local-exec" {
    command = <<-EOT
      # Login to DigitalOcean Container Registry
      echo "${var.token}" | docker login ${var.registry} -u token --password-stdin

      # Build the image
      docker build \
        -t ${local.tag} \
        -f ${var.context}/${local.dockerfile} \
        ${local.build_args} \
        ${local.target_arg} \
        ${var.context}

      # Push to registry
      docker push ${local.tag}
    EOT
  }
}
