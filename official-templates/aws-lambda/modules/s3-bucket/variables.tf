variable "name" {
  description = "Bucket name prefix"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "versioning" {
  description = "Enable versioning"
  type        = bool
  default     = false
}

variable "public" {
  description = "Allow public access"
  type        = bool
  default     = false
}
