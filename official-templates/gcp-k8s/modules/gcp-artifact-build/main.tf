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
  image_tag  = "${var.registry}/${random_id.tag.hex}"
  build_args = join(" ", [for k, v in var.args : "--build-arg ${k}=${v}"])
  target_arg = var.target != null ? "--target ${var.target}" : ""
}

resource "random_id" "tag" {
  byte_length = 8

  keepers = {
    context    = var.context
    dockerfile = local.dockerfile
    target     = var.target != null ? var.target : ""
  }
}

resource "null_resource" "docker_build" {
  triggers = {
    tag = random_id.tag.hex
  }

  provisioner "local-exec" {
    command = <<-EOT
      docker build \
        -t ${local.image_tag} \
        -f ${var.context}/${local.dockerfile} \
        ${local.build_args} \
        ${local.target_arg} \
        ${var.context} && \
      docker push ${local.image_tag}
    EOT
  }
}
