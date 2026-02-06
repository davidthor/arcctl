variable "name" {
  description = "Observability configuration name"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "log_group" {
  description = "CloudWatch log group name"
  type        = string
}
