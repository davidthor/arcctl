variable "name" {
  description = "Database instance name"
  type        = string
}

variable "type" {
  description = "Database type (postgres, mysql, redis)"
  type        = string

  validation {
    condition     = contains(["postgres", "mysql", "redis"], var.type)
    error_message = "Database type must be one of: postgres, mysql, redis."
  }
}

variable "engine_version" {
  description = "Database engine version"
  type        = string
  default     = null
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 20
}

variable "security_group_id" {
  description = "Security group ID for database access"
  type        = string
}
