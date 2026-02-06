variable "name" {
  description = "Key pair name (used as Secrets Manager path prefix)"
  type        = string
}

variable "key_type" {
  description = "Key type (rsa or ecdsa)"
  type        = string

  validation {
    condition     = contains(["rsa", "ecdsa"], var.key_type)
    error_message = "Key type must be 'rsa' or 'ecdsa'."
  }
}

variable "key_size" {
  description = "Key size in bits (RSA: 2048/4096, ECDSA: 256/384/521)"
  type        = number
  default     = null
}

variable "region" {
  description = "AWS region"
  type        = string
}
