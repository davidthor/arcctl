terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

# Resend provides SMTP access using a fixed configuration.
# The API key serves as the SMTP password - no resource creation needed.
# Host: smtp.resend.com, Port: 465, Username: resend, Password: <api_key>
resource "null_resource" "resend_smtp" {
  triggers = {
    name    = var.name
    api_key = var.api_key
  }
}
