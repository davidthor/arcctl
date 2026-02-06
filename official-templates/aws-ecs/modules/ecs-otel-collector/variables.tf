variable "name" {
  description = "OTel collector name"
  type        = string
}

variable "cluster" {
  description = "ECS cluster name"
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

variable "security_group_id" {
  description = "Security group ID"
  type        = string
}

variable "log_group" {
  description = "CloudWatch log group name"
  type        = string
}
