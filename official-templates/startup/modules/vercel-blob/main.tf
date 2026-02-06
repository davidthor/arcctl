terraform {
  required_providers {
    vercel = {
      source  = "vercel/vercel"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "vercel" {
  api_token = var.token
  team      = var.team_id != "" ? var.team_id : null
}

# Generate a unique store name
resource "random_id" "store" {
  byte_length = 4
  prefix      = "${var.name}-"
}

# Vercel Blob store is provisioned via the Vercel API.
# The store token is stored as a project environment variable
# that the Vercel Blob SDK uses automatically.
resource "null_resource" "blob_store" {
  triggers = {
    name = random_id.store.hex
  }

  provisioner "local-exec" {
    command = <<-EOT
      curl -s -X POST "https://api.vercel.com/v1/blob/stores" \
        -H "Authorization: Bearer ${var.token}" \
        -H "Content-Type: application/json" \
        ${var.team_id != "" ? "-H \"x-vercel-team: ${var.team_id}\"" : ""} \
        -d '{"name": "${random_id.store.hex}", "public": ${var.public ? "true" : "false"}}' \
        > /tmp/vercel-blob-${random_id.store.hex}.json
    EOT
  }
}

# Read the blob store details after creation
data "local_file" "blob_response" {
  filename   = "/tmp/vercel-blob-${random_id.store.hex}.json"
  depends_on = [null_resource.blob_store]
}

locals {
  blob_data = jsondecode(data.local_file.blob_response.content)
}
