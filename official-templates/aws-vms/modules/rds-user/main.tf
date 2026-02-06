terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Generate a random password for the new database user
resource "random_password" "this" {
  length  = 32
  special = false
}

# Note: The actual database user creation requires a database-specific provider
# (e.g., postgresql, mysql). Since the database type is determined at runtime,
# this module generates credentials that the application can use.
# In production, consider using a provisioner or external data source
# to create the user directly on the database.

locals {
  # Parse database connection info from the database URL
  # Expected format: protocol://username:password@host:port/database
  connection_url = var.database
  username       = var.username
  password       = random_password.this.result
}
