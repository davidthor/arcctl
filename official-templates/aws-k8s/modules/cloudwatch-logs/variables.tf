variable "name" {
  description = "CloudWatch log group name"
  type        = string
}

variable "retention_days" {
  description = "Number of days to retain log events"
  type        = number
  default     = 30
}
