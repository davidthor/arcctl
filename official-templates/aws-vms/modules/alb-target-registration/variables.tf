variable "name" {
  description = "Service name"
  type        = string
}

variable "target_group_arn" {
  description = "ALB target group ARN"
  type        = string
}

variable "target" {
  description = "Target instance ID or IP"
  type        = string
}

variable "target_type" {
  description = "Target type"
  type        = string
  default     = "instance"
}

variable "port" {
  description = "Target port"
  type        = number
}
