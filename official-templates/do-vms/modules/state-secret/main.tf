terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# State-based secret storage - secrets are stored directly in OpenTofu state.
# This is simple but means secrets are only as secure as the state backend.

resource "random_id" "secret_id" {
  byte_length = 8
  prefix      = "${var.name}-"
}

# The actual secret data is stored in state via the variable.
# No external resources are created - this is intentional for VM-based
# deployments where there's no secret manager (like K8s Secrets or Vault).
