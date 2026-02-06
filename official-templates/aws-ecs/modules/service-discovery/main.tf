terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_vpc" "this" {
  filter {
    name   = "tag:Name"
    values = ["*"]
  }
}

# Create or look up the namespace
resource "aws_service_discovery_private_dns_namespace" "this" {
  name        = "${var.namespace}.local"
  description = "Service discovery namespace for ${var.namespace}"
  vpc         = data.aws_vpc.this.id

  tags = {
    Name      = var.namespace
    ManagedBy = "arcctl"
  }
}

resource "aws_service_discovery_service" "this" {
  name = var.name

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.this.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}
