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

# Generate a unique tag for this build
resource "random_id" "tag" {
  byte_length = 8
}

locals {
  # Vercel uses their own registry for container deployments
  registry   = "registry.vercel.com"
  image_tag  = "${local.registry}/${var.team_id}/${random_id.tag.hex}"
  dockerfile = var.dockerfile != null ? var.dockerfile : "Dockerfile"
  target     = var.target != null ? "--target ${var.target}" : ""

  # Build --build-arg flags from the args map
  build_args = var.args != null ? join(" ", [
    for k, v in var.args : "--build-arg ${k}=${v}"
  ]) : ""
}

# Build and push the Docker image
resource "null_resource" "docker_build" {
  triggers = {
    context    = var.context
    dockerfile = local.dockerfile
    target     = var.target
    tag        = random_id.tag.hex
  }

  provisioner "local-exec" {
    command = <<-EOT
      # Log in to the Vercel container registry
      echo "${var.token}" | docker login ${local.registry} -u _ --password-stdin

      # Build the image
      docker build \
        -t ${local.image_tag} \
        -f ${var.context}/${local.dockerfile} \
        ${local.target} \
        ${local.build_args} \
        ${var.context}

      # Push the image
      docker push ${local.image_tag}
    EOT
  }
}
