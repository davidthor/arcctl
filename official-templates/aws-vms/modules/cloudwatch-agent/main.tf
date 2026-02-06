terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

# The CloudWatch agent is installed on each EC2 instance via user_data
# in the ec2-docker and ec2-runtime modules. This module provides
# the configuration and outputs needed for observability integration.
#
# The OTel collector sidecar runs on each instance, receiving OTLP
# data and exporting to CloudWatch Logs, Metrics, and X-Ray.

# SSM Parameter for OTel collector config
resource "aws_ssm_parameter" "otel_config" {
  name  = "/arcctl/${var.name}/otel-config"
  type  = "String"
  value = yamlencode({
    receivers = {
      otlp = {
        protocols = {
          grpc = { endpoint = "0.0.0.0:4317" }
          http = { endpoint = "0.0.0.0:4318" }
        }
      }
    }
    exporters = {
      awsxray = {
        region = data.aws_region.current.name
      }
      awsemf = {
        log_group_name = var.log_group
        region         = data.aws_region.current.name
      }
    }
    service = {
      pipelines = {
        traces = {
          receivers = ["otlp"]
          exporters = ["awsxray"]
        }
        metrics = {
          receivers = ["otlp"]
          exporters = ["awsemf"]
        }
      }
    }
  })

  tags = {
    Name      = "${var.name}-otel-config"
    ManagedBy = "arcctl"
  }
}
