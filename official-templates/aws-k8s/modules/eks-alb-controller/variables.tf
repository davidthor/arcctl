variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "kubeconfig" {
  description = "Kubeconfig for cluster access"
  type        = string
  sensitive   = true
}
